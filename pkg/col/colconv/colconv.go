// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colconv

import (
	"bytes"
	"context"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/colengine"
	"github.com/cockroachdb/cockroach/pkg/col/colserde"
	"github.com/cockroachdb/cockroach/pkg/col/coltypes"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage/engine"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

// SSTableToColumnar converts the kv data in an sstable into a columnar batch.
func SSTableToColumnar(
	ctx context.Context, schema colengine.SchemaProvider, span roachpb.Span, sstable []byte,
) (retBuf []byte, retSchema *colengine.IndexSchema, retErr error) {
	iter, err := engine.NewMemSSTIterator(sstable, false)
	if err != nil {
		return nil, nil, err
	}
	defer iter.Close()

	var prevSchema *colengine.IndexSchema
	var b coldata.Batch
	var buf bytes.Buffer
	var s *colserde.FileSerializer
	var rowIdx int
	flushBatch := func() error {
		b.SetLength(uint16(rowIdx))

		if s == nil {
			var err error
			if s, err = colserde.NewFileSerializer(&buf, prevSchema.Types); err != nil {
				return err
			}
		}
		if err := s.AppendBatch(b); err != nil {
			return err
		}

		rowIdx = 0
		b.Reset(prevSchema.Types, int(coldata.BatchSize()))
		return nil
	}

	for iter.Seek(engine.MakeMVCCMetadataKey(span.Key)); ; iter.Next() {
		ok, err := iter.Valid()
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			break
		}
		unsafeKey := iter.UnsafeKey()
		schema, err := schema.Schema(ctx, unsafeKey.Key, unsafeKey.Timestamp)
		if err != nil {
			return nil, nil, err
		}

		if prevSchema == nil {
			prevSchema = schema
			b = coldata.NewMemBatch(schema.Types)
			log.Info(ctx, "set schema ", prevSchema)
		} else if prevSchema != schema {
			return nil, nil, errors.New(`unimplemented: multiple schemas in one sstable`)
		}

		if rowIdx == int(coldata.BatchSize()) {
			if err := flushBatch(); err != nil {
				return nil, nil, err
			}
		}

		key, value := []byte(iter.UnsafeKey().Key), iter.UnsafeValue()
		if err := decodeRow(schema, b, key, value, rowIdx); err != nil {
			return nil, nil, err
		}
		rowIdx++
	}
	if prevSchema == nil {
		return nil, nil, nil
	}
	if err := flushBatch(); err != nil {
		return nil, nil, err
	}
	if err := s.Finish(); err != nil {
		return nil, nil, err
	}

	return buf.Bytes(), prevSchema, nil
}

func decodeRow(
	schema *colengine.IndexSchema, b coldata.Batch, key, value []byte, rowIdx int,
) error {
	// Skip table id and index id.
	for i := 0; i < 2; i++ {
		l, err := encoding.PeekLength(key)
		if err != nil {
			return err
		}
		key = key[l:]
	}
	for _, colID := range schema.KeyColIDs {
		rest, err := decodeKeyCol(schema, b, key, rowIdx, colID)
		if err != nil {
			return err
		}
		key = rest
	}
	v := roachpb.Value{RawBytes: value}
	valueBytes, err := v.GetTuple()
	if err != nil {
		return err
	}
	var colID uint32
	for len(valueBytes) > 0 {
		_, _, colIDDiff, typ, err := encoding.DecodeValueTag(valueBytes)
		if err != nil {
			return err
		}
		colID += colIDDiff
		rest, err := decodeValueCol(schema, b, valueBytes, rowIdx, colID, typ)
		if err != nil {
			return err
		}
		valueBytes = rest
	}
	return nil
}

func decodeKeyCol(
	schema *colengine.IndexSchema, b coldata.Batch, key []byte, rowIdx int, colID uint32,
) ([]byte, error) {
	colIdx := schema.ColIdxByColID[colID]
	colType := schema.Types[colIdx]
	switch colType {
	case coltypes.Int32, coltypes.Int64:
		rest, x, err := encoding.DecodeVarintAscending(key)
		if err != nil {
			return nil, err
		}
		b.ColVec(colIdx).Int64()[rowIdx] = x
		return rest, nil
	}
	return nil, errors.Errorf(`unhandled key type: %s`, colType)
}

func decodeValueCol(
	schema *colengine.IndexSchema,
	b coldata.Batch,
	value []byte,
	rowIdx int,
	colID uint32,
	typ encoding.Type,
) ([]byte, error) {
	colIdx, ok := schema.ColIdxByColID[colID]
	if !ok {
		return nil, errors.AssertionFailedf(`could not find column with id %d: %v`, colID, schema.ColIdxByColID)
	}
	// colTyp := schema.Types[colIdx]

	switch typ {
	case encoding.Int:
		rest, x, err := encoding.DecodeIntValue(value)
		if err != nil {
			return nil, err
		}
		b.ColVec(colIdx).Int64()[rowIdx] = x
		return rest, nil
	case encoding.Bytes:
		rest, x, err := encoding.DecodeBytesValue(value)
		if err != nil {
			return nil, err
		}
		b.ColVec(colIdx).Bytes().Set(rowIdx, x)
		return rest, nil
	case encoding.Float:
		rest, x, err := encoding.DecodeFloatValue(value)
		if err != nil {
			return nil, err
		}
		b.ColVec(colIdx).Float64()[rowIdx] = x
		return rest, nil
	}
	return nil, errors.Errorf(`unhandled value type: %s`, typ)
}
