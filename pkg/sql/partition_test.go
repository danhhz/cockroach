// Copyright 2017 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package sql_test

import (
	"bytes"
	gosql "database/sql"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/pkg/errors"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"golang.org/x/net/context"
)

func setZoneConfig(sqlDB *gosql.DB, names []string, cfg *config.ZoneConfig) error {
	save := func(objID uint32, cfg *config.ZoneConfig) error {
		buf, err := protoutil.Marshal(cfg)
		if err != nil {
			return err
		}
		_, err = sqlDB.Exec(`UPSERT INTO system.zones VALUES ($1, $2)`, objID, buf)
		return err
	}

	switch len(names) {
	case 2:
		tableID, err := sqlutils.QueryTableID(sqlDB, names[0], names[1])
		if err != nil {
			return err
		}
		return save(tableID, cfg)
	case 3:
		tableID, err := sqlutils.QueryTableID(sqlDB, names[0], names[1])
		if err != nil {
			return err
		}

		var buf []byte
		row := sqlDB.QueryRow(`SELECT config FROM system.zones WHERE ID = $1`, tableID)
		if err := row.Scan(&buf); err != nil && err != gosql.ErrNoRows {
			return err
		}
		var tableCfg config.ZoneConfig
		if err := tableCfg.Unmarshal(buf); err != nil {
			return err
		}

		// TODO(dan): So hacky
		var descBytes []byte
		if err := sqlDB.QueryRow(`SELECT d.descriptor FROM system.descriptor d JOIN system.namespace n ON d.id = n.id WHERE n.name = $1`, names[1]).Scan(&descBytes); err != nil {
			return err
		}
		var desc sqlbase.Descriptor
		if err := desc.Unmarshal(descBytes); err != nil {
			return nil
		}
		tableDesc := desc.GetTable()
		if tableDesc == nil {
			return fmt.Errorf("unable to decode table descriptor for %s", tableDesc)
		}
		if !tableDesc.PrimaryIndex.HasPartition(names[2]) {
			return fmt.Errorf("table %s has no partition named %s", names[1], names[2])
		}

		// TODO(benesch): move this logic to partition-time instead of
		// zone-config-set time.
		spansByPartition, err := tableDesc.IndexSpansByPartition(&tableDesc.PrimaryIndex)
		if err != nil {
			return err
		}
		tableCfg.PartitionSpans = make([]config.PartitionSpan, 0, len(spansByPartition))
		for partition, spans := range spansByPartition {
			for _, span := range spans {
				tableCfg.PartitionSpans = append(tableCfg.PartitionSpans, config.PartitionSpan{
					Span: span, Partition: partition,
				})
			}
		}
		sort.Slice(tableCfg.PartitionSpans, func(i, j int) bool {
			return tableCfg.PartitionSpans[i].Span.Key.Compare(tableCfg.PartitionSpans[j].Span.Key) < 0
		})

		if tableCfg.PartitionZones == nil {
			tableCfg.PartitionZones = make(map[string]config.ZoneConfig)
		}
		tableCfg.PartitionZones[names[2]] = *cfg
		return save(tableID, &tableCfg)
	}
	return errors.Errorf("could not set config for: %v", names)
}

func allRangeDescriptors(ctx context.Context, db *client.DB) ([]roachpb.RangeDescriptor, error) {
	rows, err := db.Scan(ctx, keys.Meta2Prefix, keys.MetaMax, 0)
	if err != nil {
		// NB: Don't wrap this error, as wrapped HandledRetryableTxnErrors are not
		// automatically retried by db.Txn.
		//
		// TODO(benesch): teach the KV layer to use errors.Cause.
		return nil, err
	}

	rangeDescs := make([]roachpb.RangeDescriptor, len(rows))
	for i, row := range rows {
		if err := row.ValueProto(&rangeDescs[i]); err != nil {
			return nil, errors.Wrapf(err, "%s: unable to unmarshal range descriptor", row.Key)
		}
	}
	return rangeDescs, nil
}

func TestPartitioning(t *testing.T) {
	defer leaktest.AfterTest(t)()

	cfg := config.DefaultZoneConfig()
	cfg.NumReplicas = 1
	defer config.TestingSetDefaultZoneConfig(cfg)()

	ctx := context.Background()
	tcArgs := base.TestClusterArgs{ServerArgsPerNode: map[int]base.TestServerArgs{
		0: base.TestServerArgs{StoreSpecs: []base.StoreSpec{{InMemory: true, Attributes: roachpb.Attributes{Attrs: []string{"dc0"}}}}},
		1: base.TestServerArgs{StoreSpecs: []base.StoreSpec{{InMemory: true, Attributes: roachpb.Attributes{Attrs: []string{"dc1"}}}}},
		2: base.TestServerArgs{StoreSpecs: []base.StoreSpec{{InMemory: true, Attributes: roachpb.Attributes{Attrs: []string{"dc2"}}}}},
	}}
	tc := testcluster.StartTestCluster(t, 3, tcArgs)
	defer tc.Stopper().Stop(ctx)

	verifyScansOnNode := func(query string, node string) error {
		rows, err := tc.Conns[0].Query(
			fmt.Sprintf(`SELECT context, message FROM [SHOW TRACE FOR %s]`, query),
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		var scansWrongNode []string
		var traceLines []string
		var context, message gosql.NullString
		for rows.Next() {
			if err := rows.Scan(&context, &message); err != nil {
				return err
			}
			traceLine := fmt.Sprintf("%s %s", context.String, message.String)
			traceLines = append(traceLines, traceLine)
			if !strings.Contains(message.String, "Scan") || strings.Contains(message.String, "sending batch") {
				continue
			}
			if !strings.Contains(context.String, node) {
				scansWrongNode = append(scansWrongNode, traceLine)
			}
		}
		if len(scansWrongNode) > 0 {
			var err bytes.Buffer
			fmt.Fprintf(&err, "expected scans on %s:\n%s\nfull trace:", node, strings.Join(scansWrongNode, "\n"))
			for _, traceLine := range traceLines {
				err.WriteString("\n  ")
				err.WriteString(traceLine)
			}
			return errors.New(err.String())
		}
		return nil
	}

	t.Run("List", func(t *testing.T) {
		sqlDB := sqlutils.MakeSQLRunner(t, tc.Conns[0])
		sqlDB.Exec(`CREATE DATABASE IF NOT EXISTS data`)
		sqlDB.Exec(`CREATE TABLE data.list (
			a INT, b INT, PRIMARY KEY (a, b)
		) PARTITION BY LIST (a) (
			PARTITION dc1 VALUES (10),
			PARTITION dc2 VALUES (20)
		)`)

		{
			cfg := config.DefaultZoneConfig()
			cfg.Constraints = config.Constraints{Constraints: []config.Constraint{{
				Type: config.Constraint_REQUIRED, Value: "dc1",
			}}}
			if err := setZoneConfig(sqlDB.DB, []string{"data", "list", "dc1"}, &cfg); err != nil {
				t.Fatalf("%+v", err)
			}
		}

		{
			cfg := config.DefaultZoneConfig()
			cfg.Constraints = config.Constraints{Constraints: []config.Constraint{{
				Type: config.Constraint_REQUIRED, Value: "dc2",
			}}}
			if err := setZoneConfig(sqlDB.DB, []string{"data", "list", "dc2"}, &cfg); err != nil {
				t.Fatalf("%+v", err)
			}
		}

		// TODO(dan): I'd love to avoid this SucceedsSoon, maybe if we SCATTER
		// first?
		testutils.SucceedsSoon(t, func() error {
			dc1Err := verifyScansOnNode(`SELECT * FROM data.list WHERE a = 10`, `n2`)
			dc2Err := verifyScansOnNode(`SELECT * FROM data.list WHERE a = 20`, `n3`)

			if dc1Err != nil {
				log.Info(ctx, dc1Err)
			}
			if dc2Err != nil {
				log.Info(ctx, dc2Err)
			}

			if dc1Err != nil {
				return dc1Err
			}
			if dc2Err != nil {
				return dc2Err
			}
			return nil
		})
	})

	t.Run("Range", func(t *testing.T) {
		sqlDB := sqlutils.MakeSQLRunner(t, tc.Conns[0])
		sqlDB.Exec(`CREATE DATABASE IF NOT EXISTS data`)
		sqlDB.Exec(`CREATE TABLE data.range (
			a INT, b INT, PRIMARY KEY (a, b)
		) PARTITION BY RANGE (a) (
			PARTITION dc1 VALUES LESS THAN (4),
			PARTITION dc2 VALUES LESS THAN (8)
		)`)

		{
			cfg := config.DefaultZoneConfig()
			cfg.Constraints = config.Constraints{Constraints: []config.Constraint{{
				Type: config.Constraint_REQUIRED, Value: "dc1",
			}}}
			if err := setZoneConfig(sqlDB.DB, []string{"data", "range", "dc1"}, &cfg); err != nil {
				t.Fatalf("%+v", err)
			}
		}

		{
			cfg := config.DefaultZoneConfig()
			cfg.Constraints = config.Constraints{Constraints: []config.Constraint{{
				Type: config.Constraint_REQUIRED, Value: "dc2",
			}}}
			if err := setZoneConfig(sqlDB.DB, []string{"data", "range", "dc2"}, &cfg); err != nil {
				t.Fatalf("%+v", err)
			}
		}

		// TODO(dan): I'd love to avoid this SucceedsSoon, maybe if we SCATTER
		// first?
		testutils.SucceedsSoon(t, func() error {
			rs := sqlDB.QueryStr(`SHOW TESTING_RANGES FROM TABLE data.range`)
			for _, r := range rs {
				log.Infof(ctx, "range %s", r)
			}

			dc1Err := verifyScansOnNode(`SELECT * FROM data.range WHERE a = 3`, `n2`)
			dc2Err := verifyScansOnNode(`SELECT * FROM data.range WHERE a = 7`, `n3`)

			if dc1Err != nil {
				log.Info(ctx, dc1Err)
			}
			if dc2Err != nil {
				log.Info(ctx, dc2Err)
			}

			if dc1Err != nil {
				return dc1Err
			}
			if dc2Err != nil {
				return dc2Err
			}
			return nil
		})
	})
}
