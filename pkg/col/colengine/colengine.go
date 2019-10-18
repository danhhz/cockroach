// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colengine

import (
	"bytes"
	"context"
	"math"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/coldb"
	"github.com/cockroachdb/cockroach/pkg/col/colserde"
	"github.com/cockroachdb/cockroach/pkg/col/coltypes"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/cockroachdb/errors"
)

// SchemaProvider gives a hook into SQL table and index schemas.
type SchemaProvider interface {
	Schema(ctx context.Context, key []byte, ts hlc.Timestamp) (*IndexSchema, error)
}

// Reader is the read interface to a columnar engine's data.
type Reader interface {
	Scan(context.Context, roachpb.Span) ([]byte, error)
}

// Writer is the write interface to a columnar engine's data.
type Writer interface {
	Ingest(context.Context, coldb.NamespaceID, *IndexSchema, coldata.Batch) error
}

// Engine is the interface that wraps the core operations of a columnar store.
type Engine interface {
	Reader
	Writer
}

type engine struct {
	schema SchemaProvider

	mu struct {
		syncutil.Mutex
		db *coldb.DB

		// TODO(dan): Don't share these under the mutex.
		buf bytes.Buffer
		b   coldata.Batch
	}
}

// Open returns a new columnar data engine writing into the given path.
func Open(dir string, schema SchemaProvider) (Engine, error) {
	db, err := coldb.Open(dir)
	if err != nil {
		return nil, err
	}
	eng := &engine{
		schema: schema,
	}
	eng.mu.db = db
	return eng, nil
}

func (e *engine) Ingest(
	ctx context.Context, namespace coldb.NamespaceID, schema *IndexSchema, b coldata.Batch,
) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	n, ok := e.mu.db.Namespace(namespace)
	if ok {
		// TODO(dan): Verify that schema matches.
	} else {
		var err error
		if n, err = e.mu.db.AddNamespace(namespace, toDBSchema(schema)); err != nil {
			return err
		}
	}
	return n.Add(b)
}

func (e *engine) Scan(ctx context.Context, span roachpb.Span) (_ []byte, retErr error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	log.Infof(ctx, "ColumnarScan %v", span)
	defer func() {
		if retErr != nil {
			log.Infof(ctx, "ColumnarScan error %+v", retErr)
		}
	}()

	// WIP get this from the key bounds
	const namespaceID = 0
	n, ok := e.mu.db.Namespace(namespaceID)
	if !ok {
		// No data.
		return nil, nil
	}
	bounds, err := columnarBounds(n, span)
	if err != nil {
		return nil, err
	}
	it, err := n.Scan(bounds)
	if err != nil {
		return nil, err
	}
	e.mu.buf.Reset()
	s, err := colserde.NewFileSerializer(&e.mu.buf, it.Schema().ColTypes)
	if err != nil {
		return nil, err
	}
	if e.mu.b == nil {
		e.mu.b = coldata.NewMemBatch(it.Schema().ColTypes)
	} else {
		e.mu.b.Reset(it.Schema().ColTypes, int(coldata.BatchSize()))
	}
	for it.Next(e.mu.b) {
		if err := s.AppendBatch(e.mu.b); err != nil {
			return nil, err
		}
	}
	if err := s.Finish(); err != nil {
		return nil, err
	}
	return e.mu.buf.Bytes(), nil
}

func columnarBounds(n *coldb.Namespace, span roachpb.Span) (coldata.Batch, error) {
	schema := n.Schema()
	keyTyps := make([]coltypes.T, len(schema.SortCols))
	for i, sortCol := range schema.SortCols {
		keyTyps[i] = schema.ColTypes[sortCol]
	}

	bounds := coldata.NewMemBatchWithSize(keyTyps, 2)
	bounds.SetLength(2)

	decode := func(key []byte, rowIdx int) error {
		var tableID sqlbase.ID
		var indexID sqlbase.IndexID
		var err error
		if key, tableID, indexID, err = sqlbase.DecodeTableIDIndexID(key); err != nil {
			return err
		}
		_, _ = tableID, indexID
		for keyIdx, typ := range keyTyps {
			colIdx := int(schema.SortCols[keyIdx])
			switch typ {
			case coltypes.Int64:
				if len(key) == 0 {
					if rowIdx == 0 {
						bounds.ColVec(colIdx).Int64()[rowIdx] = math.MinInt64
					} else {
						bounds.ColVec(colIdx).Int64()[rowIdx] = math.MaxInt64
					}
					continue
				}
				var x int64
				key, x, err = encoding.DecodeVarintAscending(key)
				if err != nil {
					return err
				}
				bounds.ColVec(colIdx).Int64()[rowIdx] = x
			default:
				return errors.Errorf(`unhandled key type: %s`, typ)
			}
		}
		return nil
	}

	if err := decode(span.Key, 0); err != nil {
		return nil, err
	}
	if err := decode(span.EndKey, 1); err != nil {
		return nil, err
	}

	return bounds, nil
}

// IngestAll decodes all batches in a file and ingests them.
func IngestAll(
	ctx context.Context, w Writer, namespace coldb.NamespaceID, schema *IndexSchema, data []byte,
) error {
	d, err := colserde.NewFileDeserializerFromBytes(data)
	if err != nil {
		return err
	}
	defer d.Close()

	b := coldata.NewMemBatch(d.Typs())
	for batchIdx := 0; batchIdx < d.NumBatches(); batchIdx++ {
		if err := d.GetBatch(batchIdx, b); err != nil {
			return err
		}
		if err := w.Ingest(ctx, namespace, schema, b); err != nil {
			return err
		}
	}
	return nil
}

// toDBSchema converts our sql-level schema into the low level one used by the
// columnar disk format.
func toDBSchema(schema *IndexSchema) coldb.Schema {
	sortCols := make([]int, len(schema.KeyColIDs))
	for sortColIdx, colID := range schema.KeyColIDs {
		sortCols[sortColIdx] = schema.ColIdxByColID[colID]
	}
	return coldb.Schema{
		ColTypes: schema.Types,
		SortCols: sortCols,
	}
}
