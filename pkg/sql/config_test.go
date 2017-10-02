// Copyright 2015 The Cockroach Authors.
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
	"fmt"
	"sort"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

var configID = sqlbase.ID(1)
var configDescKey = sqlbase.MakeDescMetadataKey(keys.MaxReservedDescID)

// forceNewConfig forces a system config update by writing a bogus descriptor with an
// incremented value inside. It then repeatedly fetches the gossip config until the
// just-written descriptor is found.
func forceNewConfig(t *testing.T, s *server.TestServer) config.SystemConfig {
	configID++
	configDesc := &sqlbase.Descriptor{
		Union: &sqlbase.Descriptor_Database{
			Database: &sqlbase.DatabaseDescriptor{
				Name:       "sentinel",
				ID:         configID,
				Privileges: &sqlbase.PrivilegeDescriptor{},
			},
		},
	}

	// This needs to be done in a transaction with the system trigger set.
	if err := s.DB().Txn(context.TODO(), func(ctx context.Context, txn *client.Txn) error {
		if err := txn.SetSystemConfigTrigger(); err != nil {
			return err
		}
		return txn.Put(ctx, configDescKey, configDesc)
	}); err != nil {
		t.Fatal(err)
	}
	return waitForConfigChange(t, s)
}

func waitForConfigChange(t *testing.T, s *server.TestServer) config.SystemConfig {
	var foundDesc sqlbase.Descriptor
	var cfg config.SystemConfig
	testutils.SucceedsSoon(t, func() error {
		var ok bool
		if cfg, ok = s.Gossip().GetSystemConfig(); ok {
			if val := cfg.GetValue(configDescKey); val != nil {
				if err := val.GetProto(&foundDesc); err != nil {
					t.Fatal(err)
				}
				if id := foundDesc.GetDatabase().GetID(); id != configID {
					return errors.Errorf("expected database id %d; got %d", configID, id)
				}
				return nil
			}
		}
		return errors.Errorf("got nil system config")
	})
	return cfg
}

// TestGetZoneConfig exercises config.GetZoneConfig and the sql hook for it.
func TestGetZoneConfig(t *testing.T) {
	defer leaktest.AfterTest(t)()
	params, _ := createTestServerParams()
	srv, sqlDB, _ := serverutils.StartServer(t, params)
	defer srv.Stopper().Stop(context.TODO())
	s := srv.(*server.TestServer)

	expectedCounter := uint32(keys.MaxReservedDescID)

	defaultZoneConfig := config.DefaultZoneConfig()
	defaultZoneConfig.RangeMinBytes = 1 << 20
	defaultZoneConfig.RangeMaxBytes = 1 << 20
	defaultZoneConfig.GC.TTLSeconds = 60

	{
		buf, err := protoutil.Marshal(&defaultZoneConfig)
		if err != nil {
			t.Fatal(err)
		}
		objID := keys.RootNamespaceID
		if _, err = sqlDB.Exec(`UPDATE system.zones SET config = $2 WHERE id = $1`, objID, buf); err != nil {
			t.Fatalf("problem writing zone %+v: %s", defaultZoneConfig, err)
		}
	}

	// Naming scheme for database and tables:
	// db1 has tables tb11 and tb12
	// db2 has tables tb21 and tb22

	expectedCounter++
	db1 := expectedCounter
	if _, err := sqlDB.Exec(`CREATE DATABASE db1`); err != nil {
		t.Fatal(err)
	}

	expectedCounter++
	db2 := expectedCounter
	if _, err := sqlDB.Exec(`CREATE DATABASE db2`); err != nil {
		t.Fatal(err)
	}

	expectedCounter++
	tb11 := expectedCounter
	if _, err := sqlDB.Exec(`CREATE TABLE db1.tb1 (k INT PRIMARY KEY, v INT)`); err != nil {
		t.Fatal(err)
	}

	expectedCounter++
	tb12 := expectedCounter
	if _, err := sqlDB.Exec(`CREATE TABLE db1.tb2 (k INT PRIMARY KEY, v INT)`); err != nil {
		t.Fatal(err)
	}

	expectedCounter++
	tb21 := expectedCounter
	if _, err := sqlDB.Exec(`CREATE TABLE db2.tb1 (k INT PRIMARY KEY, v INT)`); err != nil {
		t.Fatal(err)
	}

	expectedCounter++
	if _, err := sqlDB.Exec(`CREATE TABLE db2.tb2 (k INT PRIMARY KEY, v INT)`); err != nil {
		t.Fatal(err)
	}

	expectedCounter++
	tb22 := expectedCounter
	if _, err := sqlDB.Exec(`TRUNCATE TABLE db2.tb2`); err != nil {
		t.Fatal(err)
	}

	{
		cfg := forceNewConfig(t, s)

		// We have no custom zone configs.
		testCases := []struct {
			key     roachpb.RKey
			zoneCfg config.ZoneConfig
		}{
			{roachpb.RKeyMin, defaultZoneConfig},
			{keys.MakeTablePrefix(0), defaultZoneConfig},
			{keys.MakeTablePrefix(1), defaultZoneConfig},
			{keys.MakeTablePrefix(keys.MaxReservedDescID), defaultZoneConfig},
			{keys.MakeTablePrefix(db1), defaultZoneConfig},
			{keys.MakeTablePrefix(db2), defaultZoneConfig},
			{keys.MakeTablePrefix(tb11), defaultZoneConfig},
			{keys.MakeTablePrefix(tb12), defaultZoneConfig},
			{keys.MakeTablePrefix(tb21), defaultZoneConfig},
			{keys.MakeTablePrefix(tb22), defaultZoneConfig},
		}

		for tcNum, tc := range testCases {
			zoneCfg, err := cfg.GetZoneConfigForKey(tc.key)
			if err != nil {
				t.Fatalf("#%d: err=%s", tcNum, err)
			}

			if !proto.Equal(&zoneCfg, &tc.zoneCfg) {
				t.Errorf("#%d: bad zone config.\nexpected: %+v\ngot: %+v", tcNum, tc.zoneCfg, zoneCfg)
			}
		}
	}

	// Now set some zone configs. We don't have a nice way of using table
	// names for this, so we do raw puts.
	// Here is the list of dbs/tables and whether they have a custom zone config:
	// db1: true
	//   tb1: true
	//   tb2: false
	// db1: false
	//   tb1: true
	//   tb2: false
	db1Cfg := config.ZoneConfig{
		NumReplicas: 1,
		Constraints: config.Constraints{Constraints: []config.Constraint{{Value: "db1"}}},
	}
	tb11Cfg := config.ZoneConfig{
		NumReplicas: 1,
		Constraints: config.Constraints{Constraints: []config.Constraint{{Value: "db1.tb1"}}},
	}
	tb21Cfg := config.ZoneConfig{
		NumReplicas: 1,
		Constraints: config.Constraints{Constraints: []config.Constraint{{Value: "db2.tb1"}}},
	}
	for objID, objZone := range map[uint32]config.ZoneConfig{
		db1:  db1Cfg,
		tb11: tb11Cfg,
		tb21: tb21Cfg,
	} {
		buf, err := protoutil.Marshal(&objZone)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = sqlDB.Exec(`INSERT INTO system.zones VALUES ($1, $2)`, objID, buf); err != nil {
			t.Fatalf("problem writing zone %+v: %s", objZone, err)
		}
	}

	{
		cfg := forceNewConfig(t, s)

		testCases := []struct {
			key     roachpb.RKey
			zoneCfg config.ZoneConfig
		}{
			{roachpb.RKeyMin, defaultZoneConfig},
			{keys.MakeTablePrefix(0), defaultZoneConfig},
			{keys.MakeTablePrefix(1), defaultZoneConfig},
			{keys.MakeTablePrefix(keys.MaxReservedDescID), defaultZoneConfig},
			{keys.MakeTablePrefix(db1), db1Cfg},
			{keys.MakeTablePrefix(db2), defaultZoneConfig},
			{keys.MakeTablePrefix(tb11), tb11Cfg},
			{keys.MakeTablePrefix(tb12), db1Cfg},
			{keys.MakeTablePrefix(tb21), tb21Cfg},
			{keys.MakeTablePrefix(tb22), defaultZoneConfig},
		}

		for tcNum, tc := range testCases {
			zoneCfg, err := cfg.GetZoneConfigForKey(tc.key)
			if err != nil {
				t.Fatalf("#%d: err=%s", tcNum, err)
			}

			if !proto.Equal(&zoneCfg, &tc.zoneCfg) {
				t.Errorf("#%d: bad zone config.\nexpected: %+v\ngot: %+v", tcNum, tc.zoneCfg, zoneCfg)
			}
		}
	}
}

// TODO(dan): BEFORE MERGE this should all be in TestGetZoneConfig instead
func TestPartitions(t *testing.T) {
	zones := []struct {
		id  uint32
		cfg config.ZoneConfig
	}{
		{50, config.ZoneConfig{NumReplicas: 1}},
		{51, config.ZoneConfig{
			NumReplicas: 2,
			PartitionSpans: []config.PartitionSpan{
				{Partition: "[1,1)", Span: roachpb.Span{Key: []byte{1}}},
				{Partition: "[3,4)", Span: roachpb.Span{Key: []byte{3}, EndKey: []byte{5}}},
			},
			PartitionZones: map[string]config.ZoneConfig{
				"[1,1)": {NumReplicas: 3},
				"[3,4)": {NumReplicas: 4},
			},
		}},
		{52, config.ZoneConfig{
			NumReplicas: 0,
			PartitionSpans: []config.PartitionSpan{
				{Partition: "[1,1)", Span: roachpb.Span{Key: []byte{1}}},
				{Partition: "[3,4)", Span: roachpb.Span{Key: []byte{3}, EndKey: []byte{5}}},
			},
			PartitionZones: map[string]config.ZoneConfig{
				"[1,1)": {NumReplicas: 5},
				"[3,4)": {NumReplicas: 6},
			},
		}},
	}
	descs := []struct {
		id   sqlbase.ID
		desc *sqlbase.Descriptor
	}{
		{50, sqlbase.WrapDescriptor(&sqlbase.DatabaseDescriptor{ID: 50})},
		{51, sqlbase.WrapDescriptor(&sqlbase.TableDescriptor{ID: 51, ParentID: 50})},
		{52, sqlbase.WrapDescriptor(&sqlbase.TableDescriptor{ID: 52, ParentID: 50})},
	}
	var cfg config.SystemConfig
	for _, zone := range zones {
		var v roachpb.Value
		if err := v.SetProto(&zone.cfg); err != nil {
			t.Fatal(err)
		}
		cfg.Values = append(cfg.Values, roachpb.KeyValue{
			Key:   sqlbase.MakeZoneKey(sqlbase.ID(zone.id)),
			Value: v,
		})
	}
	for _, desc := range descs {
		var v roachpb.Value
		if err := v.SetProto(desc.desc); err != nil {
			t.Fatal(err)
		}
		cfg.Values = append(cfg.Values, roachpb.KeyValue{
			Key:   sqlbase.MakeDescMetadataKey(desc.id),
			Value: v,
		})
	}
	sort.Sort(roachpb.KeyValueByKey(cfg.Values))

	t.Run("GetZoneConfigForKey", func(t *testing.T) {
		tests := []struct {
			key                 []byte
			expectedNumReplicas int32
		}{
			{append(keys.MakeTablePrefix(51), 0), 2},
			{append(keys.MakeTablePrefix(51), 1), 3},
			{append(keys.MakeTablePrefix(51), 2), 2},
			{append(keys.MakeTablePrefix(51), 3), 4},
			{append(keys.MakeTablePrefix(51), 4), 4},
			{append(keys.MakeTablePrefix(51), 5), 2},
			{append(keys.MakeTablePrefix(52), 0), 1},
			{append(keys.MakeTablePrefix(52), 1), 5},
			{append(keys.MakeTablePrefix(52), 2), 1},
			{append(keys.MakeTablePrefix(52), 3), 6},
			{append(keys.MakeTablePrefix(52), 4), 6},
			{append(keys.MakeTablePrefix(52), 5), 1},
		}
		for _, test := range tests {
			t.Run("", func(t *testing.T) {
				zone, err := cfg.GetZoneConfigForKey(test.key)
				if err != nil {
					t.Fatal(err)
				}
				if zone.NumReplicas != test.expectedNumReplicas {
					t.Errorf("got %d expected %d", zone.NumReplicas, test.expectedNumReplicas)
				}
			})
		}
	})
	t.Run("ComputeSplitKey", func(t *testing.T) {
		tests := []struct {
			key           []byte
			expectedSplit []byte
		}{
			{
				key:           keys.MakeTablePrefix(51),
				expectedSplit: keys.MakeTablePrefix(51),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 0),
				expectedSplit: append(keys.MakeTablePrefix(51), 1),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 1),
				expectedSplit: append(keys.MakeTablePrefix(51), 1),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 1, 0),
				expectedSplit: append(keys.MakeTablePrefix(51), 2),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 2),
				expectedSplit: append(keys.MakeTablePrefix(51), 2),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 3),
				expectedSplit: append(keys.MakeTablePrefix(51), 3),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 3, 0),
				expectedSplit: append(keys.MakeTablePrefix(51), 5),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 4),
				expectedSplit: append(keys.MakeTablePrefix(51), 5),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 5),
				expectedSplit: append(keys.MakeTablePrefix(51), 5),
			},
			{
				key:           append(keys.MakeTablePrefix(51), 5, 0),
				expectedSplit: keys.MakeTablePrefix(52),
			},
		}
		for _, test := range tests {
			t.Run(fmt.Sprintf("%x", test.key), func(t *testing.T) {
				split := cfg.ComputeSplitKey(test.key, roachpb.RKeyMax)
				if !bytes.Equal(split, test.expectedSplit) {
					t.Errorf("got %x expected %x", []byte(split), test.expectedSplit)
				}
			})
		}
	})
}
