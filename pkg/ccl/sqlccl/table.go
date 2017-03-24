// Copyright 2017 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/pkg/ccl/LICENSE

package sqlccl

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/ccl/storageccl"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

type archiveVirtualTabler struct {
}

func init() {
	sql.VirtualTablerHook = &archiveVirtualTabler{}
}

var _ sql.VirtualTabler = &archiveVirtualTabler{}

func (a *archiveVirtualTabler) GetVirtualTableDesc(
	tn *parser.TableName,
) (*sqlbase.TableDescriptor, error) {
	log.Info(context.TODO(), "archiveVirtualTabler.GetVirtualTableDesc ", tn)
	if tn == nil || len(tn.ArchiveDir) == 0 {
		return nil, nil
	}
	backupDesc, err := readBackupDesc(string(tn.ArchiveDir))
	if err != nil {
		return nil, err
	}
	db := databaseFromBackupDesc(backupDesc, tn.Database())
	if db == nil {
		return nil, sqlbase.NewUndefinedDatabaseError(tn.Database())
	}
	normalizedTableName := parser.ReNormalizeName(tn.Table())
	for _, desc := range backupDesc.Descriptors {
		if table := desc.GetTable(); table != nil && table.ParentID == db.ID {
			if parser.ReNormalizeName(table.Name) == normalizedTableName {
				tableCopy := *table
				tableCopy.ArchiveReference = &sqlbase.TableDescriptor_ArchiveReference{
					Dir:      string(tn.ArchiveDir),
					ParentID: db.ID,
					ID:       table.ID,
				}
				tableCopy.ID = keys.VirtualDescriptorID - 1
				log.Infof(context.TODO(), "%+v", tableCopy)
				return &tableCopy, nil
			}
		}
	}
	return nil, sqlbase.NewUndefinedTableError(tn.Table())
}

func (a *archiveVirtualTabler) GetVirtualDatabaseDesc(
	name string,
) (*sqlbase.DatabaseDescriptor, error) {
	log.Info(context.TODO(), "archiveVirtualTabler.GetVirtualDatabaseDesc ", name)
	return nil, nil
}

func readBackupDesc(dir string) (*BackupDescriptor, error) {
	ctx := context.TODO()
	store, err := storageccl.ExportStorageFromURI(ctx, string(dir))
	if err != nil {
		return nil, err
	}
	defer store.Close()
	backupDesc, err := ReadBackupDescriptor(ctx, store)
	if err != nil {
		return nil, err
	}
	return &backupDesc, nil
}

func databaseFromBackupDesc(
	backupDesc *BackupDescriptor, database string,
) *sqlbase.DatabaseDescriptor {
	for _, desc := range backupDesc.Descriptors {
		if db := desc.GetDatabase(); db != nil {
			if parser.ReNormalizeName(db.Name) == parser.ReNormalizeName(database) {
				return db
			}
		}
	}
	return nil
}
