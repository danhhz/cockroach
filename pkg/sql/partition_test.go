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
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
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
		if tableCfg.PartitionZones == nil {
			tableCfg.PartitionZones = make(map[string]config.ZoneConfig)
		}
		tableCfg.PartitionZones[names[2]] = *cfg
		return save(tableID, &tableCfg)
	}
	return errors.Errorf("could not set config for: %v", names)
}

func TestPartitioning(t *testing.T) {
	defer leaktest.AfterTest(t)()

	ctx := context.Background()
	tcArgs := base.TestClusterArgs{ServerArgsPerNode: map[int]base.TestServerArgs{
		0: base.TestServerArgs{StoreSpecs: []base.StoreSpec{{InMemory: true, Attributes: roachpb.Attributes{Attrs: []string{"dc0"}}}}},
		1: base.TestServerArgs{StoreSpecs: []base.StoreSpec{{InMemory: true, Attributes: roachpb.Attributes{Attrs: []string{"dc1"}}}}},
		2: base.TestServerArgs{StoreSpecs: []base.StoreSpec{{InMemory: true, Attributes: roachpb.Attributes{Attrs: []string{"dc2"}}}}},
	}}
	tc := testcluster.StartTestCluster(t, 3, tcArgs)
	defer tc.Stopper().Stop(ctx)
	sqlDB := sqlutils.MakeSQLRunner(t, tc.Conns[0])

	sqlDB.Exec(`CREATE DATABASE data`)
	sqlDB.Exec(`CREATE TABLE data.foo (
		a INT, b INT, PRIMARY KEY (a, b)
	) PARTITION BY LIST (a, b) (
		PARTITION dc1 (VALUES (10, 11)),
		PARTITION dc2 (VALUES (20, 21))
	)`)

	cfg := config.DefaultZoneConfig()
	cfg.Constraints = config.Constraints{Constraints: []config.Constraint{{
		Type: config.Constraint_REQUIRED, Value: "dc1",
	}}}
	if err := setZoneConfig(sqlDB.DB, []string{"data", "foo", "dc1"}, &cfg); err != nil {
		t.Fatalf("%+v", err)
	}

	cfg.Constraints = config.Constraints{Constraints: []config.Constraint{{
		Type: config.Constraint_REQUIRED, Value: "dc2",
	}}}
	if err := setZoneConfig(sqlDB.DB, []string{"data", "foo", "dc2"}, &cfg); err != nil {
		t.Fatalf("%+v", err)
	}

	verifyScansOnNode := func(query string, node string) error {
		traceLines := sqlDB.QueryStr(
			fmt.Sprintf(`SELECT context, message FROM [SHOW TRACE FOR %s]`, query),
		)
		var scansWrongNode []string
		for _, traceLine := range traceLines {
			if !strings.Contains(traceLine[1], "Scan") || strings.Contains(traceLine[1], "sending batch") {
				continue
			}
			if !strings.HasPrefix(traceLine[0], node) {
				scansWrongNode = append(scansWrongNode, strings.Join(traceLine, " "))
			}
		}
		if len(scansWrongNode) > 0 {
			var err bytes.Buffer
			fmt.Fprintf(&err, "expected scans on %s:\n%s\nfull trace:", node, strings.Join(scansWrongNode, "\n"))
			for _, traceLine := range traceLines {
				err.WriteString("\n  ")
				err.WriteString(strings.Join(traceLine, " "))
			}
			return errors.New(err.String())
		}
		return nil
	}

	testutils.SucceedsSoon(t, func() error {
		dc1Err := verifyScansOnNode(`SELECT * FROM data.foo WHERE a = 10`, `n2`)
		dc2Err := verifyScansOnNode(`SELECT * FROM data.foo WHERE a = 20`, `n3`)
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
}
