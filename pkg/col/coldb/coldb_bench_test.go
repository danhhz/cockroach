// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package coldb_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/coldb"
	"github.com/cockroachdb/cockroach/pkg/col/coltypes"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/bufalloc"
	"github.com/cockroachdb/cockroach/pkg/workload"
	"github.com/cockroachdb/cockroach/pkg/workload/bank"
	"github.com/stretchr/testify/require"
)

type workloadLoader struct {
	namespaceID coldb.NamespaceID
	batchTyps   []coltypes.T
	a           bufalloc.ByteAllocator
	batch       coldata.Batch
	bytesLoaded int64
}

func (l *workloadLoader) Load(t testing.TB, db *coldb.DB, table workload.Table) coldb.NamespaceID {
	t.Helper()

	l.namespaceID++
	if l.batch == nil {
		l.batch = coldata.NewMemBatchWithSize(nil, 0)
	}

	var n *coldb.Namespace
	for batchIdx := 0; batchIdx < table.InitialRows.NumBatches; batchIdx++ {
		l.a = l.a[:0]
		table.InitialRows.FillBatch(batchIdx, l.batch, &l.a)
		if n == nil {
			l.batchTyps = l.batchTyps[:0]
			for _, col := range l.batch.ColVecs() {
				l.batchTyps = append(l.batchTyps, col.Type())
			}
			schema := coldb.Schema{
				ColTypes: l.batchTyps,
				SortCols: sortColsForTable(table),
			}
			var err error
			n, err = db.AddNamespace(l.namespaceID, schema)
			require.NoError(t, err)
		}
		require.NoError(t, n.Add(l.batch))
	}
	l.bytesLoaded += n.Size()
	return l.namespaceID
}

func BenchmarkAdd(b *testing.B) {
	b.Run("bank/rows=1000000,payload_size=8", func(b *testing.B) {
		dir, cleanup := testutils.TempDir(b)
		defer cleanup()

		const rows, payloadSize, splits = 1000000, 8, 0
		batchSize := int(coldata.BatchSize())
		gen := bank.FromConfig(rows, batchSize, payloadSize, splits)
		tables := gen.Tables()

		db, err := coldb.Open(dir)
		require.NoError(b, err)

		b.ResetTimer()
		b.StartTimer()
		var l workloadLoader
		for i := 0; i < b.N; i++ {
			for _, table := range tables {
				_ = l.Load(b, db, table)
			}
		}
		b.StopTimer()
		b.SetBytes(l.bytesLoaded / int64(b.N))
	})
}

func BenchmarkScan(b *testing.B) {
	b.Run("bank/rows=1000000,payload_size=8", func(b *testing.B) {
		dir, cleanup := testutils.TempDir(b)
		defer cleanup()

		const rows, payloadSize, splits = 1000000, 8, 0
		batchSize := int(coldata.BatchSize())
		gen := bank.FromConfig(rows, batchSize, payloadSize, splits)
		tables := gen.Tables()

		db, err := coldb.Open(dir)
		require.NoError(b, err)

		var l workloadLoader
		var namespaceIDs []coldb.NamespaceID
		for _, table := range tables {
			namespaceID := l.Load(b, db, table)
			namespaceIDs = append(namespaceIDs, namespaceID)
		}
		b.SetBytes(l.bytesLoaded)

		// TODO(dan): Figure out how to make this work for any workload.Generator.
		scanBounds := int64Bounds(0, rows)

		b.ResetTimer()
		b.StartTimer()
		scanBatch := coldata.NewMemBatchWithSize(nil, 0)
		for i := 0; i < b.N; i++ {
			for _, namespaceID := range namespaceIDs {
				n, ok := db.Namespace(namespaceID)
				require.True(b, ok)
				it, err := n.Scan(scanBounds)
				require.NoError(b, err)
				var scannedRows int
				for it.Next(scanBatch) {
					scannedRows += int(scanBatch.Length())
				}
				require.Equal(b, rows, scannedRows)
			}
		}
		b.StopTimer()
	})
}

func sortColsForTable(t workload.Table) []int {
	switch t.Name {
	case `bank`:
		return []int{0}
	default:
		panic(t)
	}
}
