// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colscan_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	_ "github.com/cockroachdb/cockroach/pkg/ccl"
	"github.com/cockroachdb/cockroach/pkg/ccl/workloadccl"
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/workload/bank"
	"github.com/stretchr/testify/require"
)

func TestColScan(t *testing.T) {
	defer leaktest.AfterTest(t)()
	ctx := context.Background()
	tc := testcluster.StartTestCluster(t, 1, base.TestClusterArgs{
		ServerArgs: base.TestServerArgs{UseDatabase: `d`},
	})
	defer tc.Stopper().Stop(ctx)
	db := sqlutils.MakeSQLRunner(tc.ServerConn(0))
	db.Exec(t, `CREATE DATABASE d`)

	gen := bank.FromRows(10)
	_, err := workloadccl.ImportFixture(ctx, tc.ServerConn(0), gen, `d`, false, 1, false, ``)
	require.NoError(t, err)

	db.Exec(t, `SET CLUSTER SETTING sql.columnar_scan.enabled = true`)
	db.Exec(t, `SET vectorize = 'experimental_always'`)
	log.Infof(ctx, "\n\n\n%v", db.QueryStr(t, `SELECT * FROM bank`))
}

const lineitemSchema = `(
	l_orderkey      INTEGER NOT NULL,
	l_partkey       INTEGER NOT NULL,
	l_suppkey       INTEGER NOT NULL,
	l_linenumber    INTEGER NOT NULL,
	l_quantity      FLOAT NOT NULL,
	l_extendedprice FLOAT NOT NULL,
	l_discount      FLOAT NOT NULL,
	l_tax           FLOAT NOT NULL,
	l_returnflag    CHAR(1) NOT NULL,
	l_linestatus    CHAR(1) NOT NULL,
	l_shipdate      DATE NOT NULL,
	l_commitdate    DATE NOT NULL,
	l_receiptdate   DATE NOT NULL,
	l_shipinstruct  CHAR(25) NOT NULL,
	l_shipmode      CHAR(10) NOT NULL,
	l_comment       VARCHAR(44) NOT NULL,
	PRIMARY KEY (l_orderkey, l_linenumber)
)`

func TestTPCH(t *testing.T) {
	defer leaktest.AfterTest(t)()
	ctx := context.Background()
	tc := testcluster.StartTestCluster(t, 1, base.TestClusterArgs{
		ServerArgs: base.TestServerArgs{
			UseDatabase:       `d`,
			ExternalIODir:     `/Users/dan/Desktop`,
			SQLMemoryPoolSize: 1 << 30,
		},
	})
	defer tc.Stopper().Stop(ctx)
	db := sqlutils.MakeSQLRunner(tc.ServerConn(0))
	db.Exec(t, `CREATE DATABASE d`)
	db.Exec(t, `SET CLUSTER SETTING kv.bulk_ingest.pk_buffer_size = $1`, 512<<20)
	db.Exec(t, `SET CLUSTER SETTING kv.bulk_ingest.max_pk_buffer_size = $1`, 512<<20)
	db.Exec(t,
		`IMPORT TABLE lineitem `+
			lineitemSchema+` CSV DATA ('nodelocal:///lineitem.csv')`+
			`WITH delimiter='|'`)
	db.Exec(t, `ALTER TABLE lineitem SPLIT AT VALUES (2, 0)`)

	db.Exec(t, `SET CLUSTER SETTING sql.columnar_scan.enabled = true`)
	db.Exec(t, `SET vectorize = 'experimental_always'`)
	log.Infof(ctx, "\n\n%v\n\n", db.QueryStr(t, `SELECT
	sum(l_extendedprice * l_discount) AS revenue
FROM
	lineitem
WHERE
	l_shipdate >= DATE '1994-01-01'
	AND l_shipdate < DATE '1994-01-01' + INTERVAL '1' YEAR
	AND l_discount BETWEEN 0.06 - 0.01 AND 0.06 + 0.01
	AND l_quantity < 24`))
}

func benchmarkColScan(b *testing.B, useColScan bool) {
	dir, cleanup := testutils.TempDir(b)
	defer cleanup()

	ctx := context.Background()
	tc := testcluster.StartTestCluster(b, 1, base.TestClusterArgs{
		ServerArgs: base.TestServerArgs{
			UseDatabase: `d`,
			StoreSpecs:  []base.StoreSpec{{InMemory: false, Path: dir}},
		},
	})
	defer tc.Stopper().Stop(ctx)
	db := sqlutils.MakeSQLRunner(tc.ServerConn(0))
	db.Exec(b, `CREATE DATABASE d`)
	db.Exec(b, `SET CLUSTER SETTING kv.bulk_ingest.max_pk_buffer_size = $1`, 1<<30)

	const numRows, payloadBytes, ranges = 1000000, 8, 0
	gen := bank.FromConfig(numRows, int(coldata.BatchSize()), payloadBytes, ranges)
	_, err := workloadccl.ImportFixture(ctx, tc.ServerConn(0), gen, `d`, false, 1, false, ``)
	require.NoError(b, err)

	db.Exec(b, `SET CLUSTER SETTING sql.columnar_scan.enabled = $1`, useColScan)
	db.Exec(b, `SET vectorize = 'experimental_always'`)

	b.ResetTimer()
	b.StartTimer()
	var totalBytes int64
	testutils.AllocProfileDiff(b, `/tmp/mem.before`, `/tmp/mem.after`, func() {
		for i := 0; i < b.N; i++ {
			var bytes int64
			db.QueryRow(b, `SELECT sum_int(id) FROM bank`).Scan(&bytes)
			totalBytes += bytes / 8
		}
	})
	b.StopTimer()
	b.SetBytes(numRows * 8)
}

func BenchmarkColScan(b *testing.B) {
	b.Run(`row`, func(b *testing.B) { benchmarkColScan(b, false) })
	b.Run(`col`, func(b *testing.B) { benchmarkColScan(b, true) })
}
