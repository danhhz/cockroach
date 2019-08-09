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
	"math/rand"
	"sort"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/coldb"
	"github.com/cockroachdb/cockroach/pkg/col/coltypes"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/stretchr/testify/require"
)

func int64Bounds(min, max int64) coldata.Batch {
	b := coldata.NewMemBatchWithSize([]coltypes.T{coltypes.Int64}, 2)
	b.SetLength(2)
	int64Col := b.ColVec(0).Int64()
	int64Col[0] = min
	int64Col[1] = max
	return b
}

func TestBasicRoundtrip(t *testing.T) {
	dir, cleanup := testutils.TempDir(t)
	defer cleanup()

	db, err := coldb.Open(dir)
	require.NoError(t, err)
	require.NotNil(t, db)

	const namespaceID = coldb.NamespaceID(0)
	schema := coldb.Schema{
		ColTypes: []coltypes.T{coltypes.Int64},
		SortCols: []int{0},
	}
	const numRows = 10
	var testdata coldata.Batch
	{
		testdata = coldata.NewMemBatch(schema.ColTypes)
		testdata.SetLength(numRows)
		int64Col := testdata.ColVec(0).Int64()
		for i := 0; i < numRows; i++ {
			int64Col[i] = int64(numRows - i)
		}
	}

	// Write some data.
	{
		d, err := db.AddNamespace(namespaceID, schema)
		require.NoError(t, err)
		require.NotNil(t, d)
		require.NoError(t, d.Add(testdata))
	}

	// Read it back.
	{
		d, ok := db.Namespace(namespaceID)
		require.True(t, ok, `expected to find data with id: %d`, namespaceID)
		require.NotNil(t, d)

		scanBounds := int64Bounds(-1, numRows+1)
		it, err := d.Scan(scanBounds)
		require.NoError(t, err)
		require.NotNil(t, it)

		// First Next should succeed with our data.
		roundtripped := coldata.NewMemBatchWithSize(schema.ColTypes, 0)
		require.True(t, it.Next(roundtripped))
		assertEqualBatches(t, testdata, roundtripped)

		// Second Next should not find any new data.
		require.False(t, it.Next(roundtripped))
	}
}

func TestRandBounds(t *testing.T) {
	rng, _ := randutil.NewPseudoRand()
	dir, cleanup := testutils.TempDir(t)
	defer cleanup()

	db, err := coldb.Open(dir)
	require.NoError(t, err)

	const namespaceID = coldb.NamespaceID(0)
	schema := coldb.Schema{
		ColTypes: []coltypes.T{coltypes.Int64},
		SortCols: []int{0},
	}

	const maxBatches, maxRowsPerBatch = 10, 10
	const minValue = -10
	var maxValue int64
	var batches [][]int64
	{
		numBatches := rng.Intn(maxBatches) + 1
		var v int64 = minValue
		for batchIdx := 0; batchIdx < numBatches; batchIdx++ {
			numRows := rng.Intn(maxRowsPerBatch) + 1
			batch := make([]int64, numRows)
			for rowIdx := 0; rowIdx < numRows; rowIdx++ {
				batch[rowIdx] = v
				maxValue = v
				// Create some values next to each other, some with gaps.
				v += rng.Int63n(2) + 1
			}
			batches = append(batches, batch)
			// Create some batches next to each other, some with gaps.
			v += rng.Int63n(2)
		}

		// Shuffle the batches to make sure insert order doesn't matter.
		rng.Shuffle(len(batches), func(i, j int) {
			batches[i], batches[j] = batches[j], batches[i]
		})

		// Make maxValue an exclusive bound.
		maxValue++
	}
	t.Logf(" batches [%d,%d): %v", minValue, maxValue, batches)

	// Load the batches into the store.
	{
		d, err := db.AddNamespace(namespaceID, schema)
		require.NoError(t, err)
		for _, rows := range batches {
			b := coldata.NewMemBatchWithSize(schema.ColTypes, len(rows))
			b.SetLength(uint16(len(rows)))
			int64Vec := b.ColVec(0).Int64()
			for rowIdx, datum := range rows {
				int64Vec[rowIdx] = datum
			}
			require.NoError(t, d.Add(b))
		}
	}

	// Verify some scans with random bounds, with extra probability of some edge
	// conditions.
	const numScans = 5
	for i := 0; i < numScans; i++ {
		scanMin, scanMax := randScanBounds(rng, minValue, maxValue)
		// Sanity check the bounds.
		if !(scanMin < scanMax) {
			scanMax = scanMin + 1
			// WIP work out the edge cases
			// t.Fatalf("invalid scan bounds [%d,%d)", scanMin, scanMax)
		}

		var expected []int64
		for _, rows := range batches {
			for _, datum := range rows {
				if datum >= scanMin && datum < scanMax {
					expected = append(expected, datum)
				}
			}
		}
		sort.Slice(expected, func(i, j int) bool { return expected[i] < expected[j] })
		t.Logf("scan [%d,%d) expected: %v", scanMin, scanMax, expected)

		n, ok := db.Namespace(namespaceID)
		require.True(t, ok)
		it, err := n.Scan(int64Bounds(scanMin, scanMax))
		require.NoError(t, err)
		var actual []int64
		b := coldata.NewMemBatchWithSize(schema.ColTypes, 0)
		for {
			ok := it.Next(b)
			if !ok {
				break
			}
			actual = append(actual, b.ColVec(0).Int64()[:int(b.Length())]...)
		}
		require.Equal(t, expected, actual)
	}
}

func randScanBounds(rng *rand.Rand, minValue, maxValue int64) (int64, int64) {
	var scanMin, scanMax int64
	switch rng.Intn(10) {
	case 0:
		scanMin = minValue
	case 1:
		scanMin = minValue - 1
	case 2:
		scanMin = minValue + 1
	default:
		// Don't have to worry about the 0 check here because there's always at
		// least one batch and at least one row in each batch.
		scanMin = minValue + rng.Int63n(maxValue-minValue)
	}
	switch rng.Intn(10) {
	case 0:
		scanMax = maxValue
	case 1:
		scanMax = maxValue + 1
	case 2:
		scanMax = maxValue - 1
	default:
		scanRange := maxValue - scanMin
		if scanRange == 0 {
			scanMax = scanMin + 1
		} else {
			scanMax = scanMin + rng.Int63n(scanRange) + 1
		}
	}
	return scanMin, scanMax
}

func TestErrors(t *testing.T) {
	dir, cleanup := testutils.TempDir(t)
	defer cleanup()

	db, err := coldb.Open(dir)
	require.NoError(t, err)

	schema := coldb.Schema{
		ColTypes: []coltypes.T{coltypes.Int64},
		SortCols: []int{0},
	}
	_, err = db.AddNamespace(coldb.NamespaceID(0), schema)
	require.NoError(t, err)

	{
		// Adding the same id again should fail.
		_, err := db.AddNamespace(coldb.NamespaceID(0), schema)
		require.EqualError(t, err, `already have a data with id: 0`)
	}
}

// TODO(dan): Deduplicate with the one in colserde.
func assertEqualBatches(t *testing.T, expected, actual coldata.Batch) {
	t.Helper()

	if actual.Selection() != nil {
		t.Fatal("violated invariant that batches have no selection vectors")
	}
	require.Equal(t, expected.Length(), actual.Length())
	require.Equal(t, expected.Width(), actual.Width())
	for colIdx := 0; colIdx < expected.Width(); colIdx++ {
		// Verify equality of ColVecs (this includes nulls). Since the coldata.Vec
		// backing array is always of coldata.BatchSize due to the scratch batch
		// that the converter keeps around, the coldata.Vec needs to be sliced to
		// the first length elements to match on length, otherwise the check will
		// fail.
		expectedVec := expected.ColVec(colIdx)
		actualVec := actual.ColVec(colIdx)
		require.Equal(t, expectedVec.Type(), actualVec.Type())
		require.Equal(
			t,
			expectedVec.Slice(expectedVec.Type(), 0, uint64(expected.Length())),
			actualVec.Slice(actualVec.Type(), 0, uint64(actual.Length())),
		)
	}
}
