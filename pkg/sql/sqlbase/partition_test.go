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

package sqlbase

import (
	gosql "database/sql"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/pkg/errors"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
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

func TestPartitions(t *testing.T) {
	defer leaktest.AfterTest(t)()
	ctx := context.Background()
	s, conn, _ := serverutils.StartServer(t, base.TestServerArgs{})
	defer s.Stopper().Stop(ctx)
	sqlDB := sqlutils.MakeSQLRunner(t, conn)

	sqlDB.Exec(`CREATE DATABASE data`)
	sqlDB.Exec(`CREATE TABLE data.foo (
		a INT, b INT, PRIMARY KEY (a, b)
	) PARTITION BY LIST (a, b) (
		PARTITION p1 (VALUES (10, 11)),
		PARTITION p2 (VALUES (20, 21))
	)`)

	cfg := config.DefaultZoneConfig()
	cfg.NumReplicas = 5
	if err := setZoneConfig(sqlDB.DB, []string{"data", "foo", "p1"}, &cfg); err != nil {
		t.Fatalf("%+v", err)
	}

	var buf []byte
	rows := sqlDB.Query(`SELECT config FROM system.zones`)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&buf); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := cfg.Unmarshal(buf); err != nil {
			t.Fatalf("%+v", err)
		}
		log.Infof(ctx, "%v", cfg)
	}
}
