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

	"github.com/cockroachdb/cockroach/pkg/col/colengine"
	"github.com/cockroachdb/cockroach/pkg/col/coltypes"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/typeconv"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/pkg/errors"
)

type schemaProvider struct {
	leaseMgr *sql.LeaseManager
	mu       struct {
		syncutil.Mutex
		prev *colengine.IndexSchema
	}
}

var _ colengine.SchemaProvider = (*schemaProvider)(nil)

// MakeSchemaProvider wraps LeaseManager to impelement the SchemaProvider
// interface.
func MakeSchemaProvider(leaseMgr *sql.LeaseManager) colengine.SchemaProvider {
	return &schemaProvider{leaseMgr: leaseMgr}
}

func (sp *schemaProvider) Schema(
	ctx context.Context, key []byte, ts hlc.Timestamp,
) (*colengine.IndexSchema, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// TODO(dan): Do we need this cacheing?
	prevSchemaMatches := sp.mu.prev != nil &&
		bytes.HasPrefix(key, sp.mu.prev.KeyPrefix) &&
		// WIP check the equals cases here
		!sp.mu.prev.MaxValid.Less(ts) &&
		!ts.Less(sp.mu.prev.MinValid)
	if prevSchemaMatches {
		return sp.mu.prev, nil
	}

	// TODO(dan): Don't hold the lock while we do this.
	rest, tableID, indexID, err := sqlbase.DecodeTableIDIndexID(key)
	if err != nil {
		return nil, err
	}
	// TODO(dan): At some point, we do actually need to keep the leases.
	tableDesc, _, err := sp.leaseMgr.Acquire(ctx, ts, tableID)
	if err != nil {
		return nil, err
	}
	// Immediately release the lease, since we only need it for the exact
	// timestamp requested.
	if err := sp.leaseMgr.Release(tableDesc); err != nil {
		return nil, err
	}
	indexDesc, err := tableDesc.FindIndexByID(indexID)
	if err != nil {
		return nil, err
	}

	schema := &colengine.IndexSchema{
		KeyPrefix:     key[:len(key)-len(rest)],
		Types:         make([]coltypes.T, len(indexDesc.ColumnIDs), len(tableDesc.Columns)),
		ColIdxByColID: make(map[uint32]int, len(indexDesc.ColumnIDs)),
		KeyColIDs:     make([]uint32, len(indexDesc.ColumnIDs)),
		MinValid:      tableDesc.ModificationTime,
		MaxValid:      ts,
	}
	for colIdx, colID := range indexDesc.ColumnIDs {
		schema.KeyColIDs[colIdx] = uint32(colID)
		schema.ColIdxByColID[uint32(colID)] = colIdx
		col, err := tableDesc.FindColumnByID(colID)
		if err != nil {
			return nil, err
		}
		schema.Types[colIdx] = typeconv.FromColumnType(&col.Type)
	}
	if indexDesc.ID == tableDesc.PrimaryIndex.ID {
		for i := range tableDesc.Columns {
			col := &tableDesc.Columns[i]
			if _, ok := schema.ColIdxByColID[uint32(col.ID)]; ok {
				// This column is already in the key.
				continue
			}
			colIdx := len(schema.Types)
			schema.Types = append(schema.Types, typeconv.FromColumnType(&col.Type))
			schema.ColIdxByColID[uint32(col.ID)] = colIdx
		}
	} else {
		return nil, errors.New(`WIP support secondary indexes in ColumnarSchema`)
	}

	sp.mu.prev = schema
	return schema, nil
}
