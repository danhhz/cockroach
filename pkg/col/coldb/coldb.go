// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package coldb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/colserde"
	"github.com/cockroachdb/cockroach/pkg/col/coltypes"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/cockroachdb/errors"
)

// TODO
// - MVCC

// NamespaceID is a unique identifer for a Namespace within a Store.
type NamespaceID uint64

// DB is a collection of durable columnar data, organized into "namespaces",
// each of which has a consistent schema.
type DB struct {
	dir string
	mu  struct {
		syncutil.RWMutex
		namespacesByID map[NamespaceID]*Namespace
	}
}

// Open returns a Store that writes data to the specified directory.
func Open(dir string) (*DB, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		return nil, errors.New(`unimplemented`)
	}
	db := &DB{dir: dir}
	db.mu.namespacesByID = make(map[NamespaceID]*Namespace)
	return db, nil
}

// Namespace returns the Namespace with the supplied ID.
func (db *DB) Namespace(id NamespaceID) (*Namespace, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	d, ok := db.mu.namespacesByID[id]
	return d, ok
}

// AddNamespace adds a Namespace with the given id and schema.
func (db *DB) AddNamespace(id NamespaceID, schema Schema) (*Namespace, error) {
	if len(schema.SortCols) == 0 {
		return nil, errors.New(`invalid sort columns`)
	}
	for colIdx := range schema.SortCols {
		if colIdx < 0 || colIdx >= len(schema.ColTypes) {
			return nil, errors.New(`invalid sort columns`)
		}
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.mu.namespacesByID[id]; ok {
		return nil, errors.Errorf(`already have a data with id: %d`, id)
	}
	n := &Namespace{
		dir:    db.dir,
		schema: schema,
	}
	db.mu.namespacesByID[id] = n
	return n, nil
}

// Iterator is an iterator over durable columnar data, all with the same schema.
type Iterator struct {
	schema  Schema
	bounds  coldata.Batch
	files   []namespaceFile
	fileIdx int
}

// Schema returns the types of the data returned by this iterator.
func (i *Iterator) Schema() Schema {
	return i.schema
}

// Next fills the next batch of data.
func (i *Iterator) Next(b coldata.Batch) bool {
	for i.fileIdx < len(i.files) {
		srcBatch := i.files[i.fileIdx].batch
		i.fileIdx++
		srcStartIdx := firstGEBatchIdx(srcBatch, i.bounds, 0, i.schema.SortCols)
		srcEndIdx := firstGEBatchIdx(srcBatch, i.bounds, 1, i.schema.SortCols)
		if srcStartIdx > srcEndIdx {
			panic(`shouldn't happen`)
		} else if srcStartIdx == srcEndIdx {
			continue
		}

		b.Reset(i.schema.ColTypes, 0)

		// WIP we currently allow execution operators to modify input data, which
		// makes this unsafe, they'll poison the source mmap'd data. but if we don't
		// allow them to do that, then we can avoid ever memcpy'ing the mmap'd data,
		// which seems like a worthwhile optimization? maybe introducing mvcc takes
		// this away from us again
		dstColVecs := b.ColVecs()
		for colIdx, colType := range i.schema.ColTypes {
			srcCol := srcBatch.ColVec(colIdx)
			dstColVecs[colIdx] = srcCol.Slice(colType, uint64(srcStartIdx), uint64(srcEndIdx))
		}
		b.SetColVecs(dstColVecs)
		b.SetLength(uint16(srcEndIdx - srcStartIdx))
		return true
	}
	return false
}

type namespaceFile struct {
	path   string
	size   int64
	bounds coldata.Batch
	batch  coldata.Batch
	des    *colserde.FileDeserializer
}

// Schema describes the format of a set of sorted columnar data.
type Schema struct {
	ColTypes []coltypes.T
	// TODO(dan): Column directions.
	SortCols []int
}

// Namespace is a collection of durable columnar data, all with the same schema
// and kept sorted by some set of the columns.
type Namespace struct {
	dir    string
	schema Schema
	mu     struct {
		syncutil.RWMutex
		size  int64
		files []namespaceFile
	}
	compactionMu struct {
		syncutil.Mutex
		filenameIdx int
	}
}

// Schema return the schema of this columnar data.
func (n *Namespace) Schema() Schema {
	return n.schema
}

// Scan returns an iterator over the given bounds.
func (n *Namespace) Scan(bounds coldata.Batch) (*Iterator, error) {
	if bounds.Width() != len(n.schema.SortCols) {
		return nil, errors.New(`invalid bounds`)
	}
	log.Infof(context.TODO(), "scanning columnar %s", DebugPrintBatch(bounds))

	_, overlapping, _ := n.overlapping(bounds)
	it := &Iterator{
		schema:  n.schema,
		bounds:  bounds,
		files:   overlapping,
		fileIdx: 0,
	}
	return it, nil
}

// Size returns the number of bytes on disk.
func (n *Namespace) Size() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.mu.size
}

// Add inserts the given batch into the data.
func (n *Namespace) Add(b coldata.Batch) error {
	cols := b.ColVecs()
	if len(n.schema.ColTypes) != len(cols) {
		return errors.New(`schemas don't match`)
	}
	for i, colType := range n.schema.ColTypes {
		if colType != cols[i].Type() {
			return errors.New(`schemas don't match`)
		}
	}
	if b.Length() == 0 {
		return nil
	}
	return n.compact(b)
}

func (n *Namespace) compact(b coldata.Batch) error {
	bounds, err := getBounds(b, n.schema.SortCols)
	if err != nil {
		return err
	}
	log.Infof(context.TODO(), "ingesting columnar %s", DebugPrintBatch(bounds))

	n.compactionMu.Lock()
	defer n.compactionMu.Unlock()
	before, overlapping, after := n.overlapping(bounds)

	if len(overlapping) > 0 {
		return errors.New(`unimplemented`)
	}

	var merged []namespaceFile
	{
		path := filepath.Join(n.dir, fmt.Sprintf("%08d.arrow", n.compactionMu.filenameIdx))
		n.compactionMu.filenameIdx++
		des, mmapBatch, size, err := writeBatchToFile(path, n.schema.ColTypes, b)
		if err != nil {
			// Clean up after ourselves.
			_ = os.Remove(path)
			return err
		}
		merged = append(merged, namespaceFile{
			path:   path,
			size:   size,
			bounds: bounds,
			batch:  mmapBatch,
			des:    des,
		})
	}

	newFiles := make([]namespaceFile, 0, len(before)+len(merged)+len(after))
	newFiles = append(newFiles, before...)
	newFiles = append(newFiles, merged...)
	newFiles = append(newFiles, after...)
	var newSize int64
	for i := range newFiles {
		newSize += newFiles[i].size
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	// TODO(dan): garbage collect the old files once no readers are using them.
	n.mu.size = newSize
	n.mu.files = newFiles
	return nil
}

// writeBatchToFile writes the given batch to the given path and returns an
// mmap'd copy of it. It's safe to reuse `b` after this returns.
func writeBatchToFile(
	path string, schema []coltypes.T, b coldata.Batch,
) (*colserde.FileDeserializer, coldata.Batch, int64, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, 0, err
	}
	ser, err := colserde.NewFileSerializer(f, schema)
	if err != nil {
		return nil, nil, 0, err
	}
	if err := ser.AppendBatch(b); err != nil {
		return nil, nil, 0, err
	}
	if err := ser.Finish(); err != nil {
		return nil, nil, 0, err
	}
	// Turn around and immediately read it back in so we can have it mmap'd and
	// the caller can reuse `b`.
	des, err := colserde.NewFileDeserializerFromPath(path)
	if err != nil {
		return nil, nil, 0, err
	}
	// TODO(dan): this leaks the deserializer and the files it opens.
	mmapBatch := coldata.NewMemBatchWithSize(schema, 0)
	if err := des.GetBatch(0, mmapBatch); err != nil {
		return nil, nil, 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, nil, 0, err
	}
	return des, mmapBatch, fi.Size(), nil
}

func getBounds(batch coldata.Batch, sortCols []int) (coldata.Batch, error) {
	l := int(batch.Length())
	if l == 0 {
		return nil, errors.Errorf(`empty batch is not supported`)
	}
	typs := make([]coltypes.T, len(sortCols))
	for i := range sortCols {
		typs[i] = batch.ColVec(int(sortCols[i])).Type()
	}
	bounds := coldata.NewMemBatchWithSize(typs, 2)
	bounds.SetLength(2)
	for i := range sortCols {
		colIdx := sortCols[i]
		boundsCol, batchCol := bounds.ColVec(colIdx), batch.ColVec(colIdx)
		firstRowIdx, lastRowIdx := uint64(0), uint64(l-1)
		boundsCol.Copy(coldata.CopySliceArgs{SliceArgs: coldata.SliceArgs{
			ColType:     typs[i],
			DestIdx:     0,
			Src:         batchCol,
			SrcStartIdx: firstRowIdx,
			SrcEndIdx:   firstRowIdx + 1,
		}})
		boundsCol.Copy(coldata.CopySliceArgs{SliceArgs: coldata.SliceArgs{
			ColType:     typs[i],
			DestIdx:     1,
			Src:         batchCol,
			SrcStartIdx: lastRowIdx,
			SrcEndIdx:   lastRowIdx + 1,
		}})
	}
	// WIP verify that the rest of the batch is sorted
	return bounds, nil
}

func (n *Namespace) overlapping(bounds coldata.Batch) (before, overlapping, after []namespaceFile) {
	// TODO(dan): binary search instead
	// WIP when testing this, make sure to check bound lengths != 2.
	n.mu.RLock()
	defer n.mu.RUnlock()

	// TODO(dan): use a binary search here, too.
	for _, file := range n.mu.files {
		fileLast := int(file.bounds.Length()) - 1
		afterEnd := compare(bounds, 0, file.bounds, fileLast, n.schema.SortCols) > 0
		if afterEnd {
			before = append(before, file)
			continue
		}
		boundsLast := int(bounds.Length()) - 1
		beforeBegin := compare(bounds, boundsLast, file.bounds, 0, n.schema.SortCols) <= 0
		if beforeBegin {
			after = append(after, file)
			continue
		}
		overlapping = append(overlapping, file)
	}
	return
}

// TODO(dan): There's got to be one of these already in exec somewhere.
func compare(x coldata.Batch, xIdx int, y coldata.Batch, yIdx int, sortCols []int) (ret int) {
	for _, sortCol := range sortCols {
		colX, colY := x.ColVec(sortCol), y.ColVec(sortCol)
		switch colX.Type() {
		case coltypes.Int64:
			valX := colX.Int64()[xIdx]
			valY := colY.Int64()[yIdx]
			if valX < valY {
				return -1
			} else if valX > valY {
				return 1
			}
		default:
			panic(errors.Errorf(`unimplemented: %s`, colX.Type()))
		}
	}
	return 0
}

func firstGEBatchIdx(
	haystack coldata.Batch, needle coldata.Batch, needleIdx int, sortCols []int,
) int {
	return sort.Search(int(haystack.Length()), func(haystackIdx int) bool {
		return compare(haystack, haystackIdx, needle, needleIdx, sortCols) >= 0
	})
}

// DebugPrintBatch WIP
func DebugPrintBatch(b coldata.Batch) string {
	var s strings.Builder
	l := int(b.Length())
	fmt.Fprintf(&s, "len=%d", l)
	if b.Width() > 1 {
		fmt.Fprintf(&s, "\n")
	}
	for colIdx, col := range b.ColVecs() {
		if colIdx > 0 {
			fmt.Fprintf(&s, "\n")
		}
		fmt.Fprintf(&s, "  %d: %s %v", colIdx, col.Type(), col.Col())
	}
	return s.String()
}
