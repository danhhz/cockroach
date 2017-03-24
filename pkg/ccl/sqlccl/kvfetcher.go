// Copyright 2017 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/pkg/ccl/LICENSE

package sqlccl

import (
	"bytes"
	"runtime/debug"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/ccl/storageccl"
	"github.com/cockroachdb/cockroach/pkg/ccl/storageccl/engineccl"
	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/storage/engine"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func init() {
	oldMakeKVFetcher := sqlbase.MakeKVFetcher
	sqlbase.MakeKVFetcher = func(desc *sqlbase.TableDescriptor) sqlbase.KVFetcher {
		if desc.ArchiveReference != nil {
			return &archiveKVFetcher{tableDesc: desc}
		}
		return oldMakeKVFetcher(desc)
	}
}

type archiveKVFetcher struct {
	tableDesc *sqlbase.TableDescriptor
	// backupDesc *BackupDescriptor
	kr storageccl.KeyRewriter

	spans   roachpb.Spans
	spanIdx int

	ssts  []engine.RocksDBSstFileReader
	iters []engine.Iterator
	iter  *engineccl.MultiIterator
}

var _ sqlbase.KVFetcher = &archiveKVFetcher{}

func (f *archiveKVFetcher) Init(
	ctx context.Context,
	_ *client.Txn,
	spans roachpb.Spans,
	reverse bool,
	_ bool,
	_ int64,
	returnRangeInfo bool,
) error {
	debug.PrintStack()
	if returnRangeInfo {
		return errors.New("range info is not supported by archiveKVFetcher")
	}
	if reverse {
		return errors.New("reverse is not supported by archiveKVFetcher")
	}

	store, err := storageccl.ExportStorageFromURI(ctx, f.tableDesc.ArchiveReference.Dir)
	defer store.Close()
	backupDesc, err := ReadBackupDescriptor(ctx, store)
	if err != nil {
		return err
	}

	var paths []string
	for _, file := range backupDesc.Files {
		localPath, err := store.FetchFile(ctx, file.Path)
		if err != nil {
			return err
		}
		paths = append(paths, localPath)
	}

	f.ssts = make([]engine.RocksDBSstFileReader, len(paths))
	f.iters = make([]engine.Iterator, len(f.ssts))
	for i, path := range paths {
		var err error
		f.ssts[i], err = engine.MakeRocksDBSstFileReader()
		if err != nil {
			return err
		}
		if err = f.ssts[i].AddFile(path); err != nil {
			return err
		}
		f.iters[i] = f.ssts[i].NewIterator(false)
	}
	f.iter = engineccl.MakeMultiIterator(f.iters)

	kr := MakeKeyRewriterForNewTableID(f.tableDesc, f.tableDesc.ArchiveReference.ID)
	for _, span := range spans {
		rewrittenStart, ok := kr.RewriteKey(span.Key)
		if !ok {
			return errors.Errorf("could not rewrite key: %s", roachpb.Key(rewrittenStart))
		}
		rewrittenEnd, ok := kr.RewriteKey(span.EndKey)
		if !ok {
			return errors.Errorf("could not rewrite key: %s", roachpb.Key(rewrittenStart))
		}
		f.spans = append(f.spans, roachpb.Span{
			Key:    rewrittenStart,
			EndKey: rewrittenEnd,
		})
	}
	f.spanIdx = -1

	tableDescCopy := *f.tableDesc
	tableDescCopy.ID = f.tableDesc.ArchiveReference.ID
	f.kr = MakeKeyRewriterForNewTableID(&tableDescCopy, f.tableDesc.ID)

	return nil
}

func (f *archiveKVFetcher) NextKV(ctx context.Context) (bool, client.KeyValue, error) {
	log.Info(ctx, f.spanIdx, " ", f.iter.Valid())
	for !f.iter.Valid() || bytes.Compare(f.iter.UnsafeKey().Key, f.spans[f.spanIdx].EndKey) >= 0 {
		log.Info(ctx, f.spanIdx, " ", f.iter.Valid())
		if err := f.iter.Error(); err != nil {
			return false, client.KeyValue{}, err
		}
		f.spanIdx++
		if f.spanIdx == len(f.spans) {
			return false, client.KeyValue{}, nil
		}
		log.Info(ctx, f.spans[f.spanIdx].Key)
		f.iter.Seek(engine.MVCCKey{Key: f.spans[f.spanIdx].Key})
	}
	unsafeKey := f.iter.UnsafeKey()
	kv := client.KeyValue{
		Key: append([]byte(nil), unsafeKey.Key...),
		Value: &roachpb.Value{
			Timestamp: unsafeKey.Timestamp,
			RawBytes:  append([]byte(nil), f.iter.UnsafeValue()...),
		},
	}
	var ok bool
	kv.Key, ok = f.kr.RewriteKey(kv.Key)
	if !ok {
		return false, client.KeyValue{}, errors.Errorf("could not rewrite key: %s", kv.Key)
	}
	log.Info(ctx, kv.Key, " -> ", kv.Value.PrettyPrint())
	f.iter.NextKey()
	return true, kv, nil
}

func (f *archiveKVFetcher) GetRangesInfo() []roachpb.RangeInfo {
	panic("GetRangeInfo() called on kvFetcher that wasn't configured with returnRangeInfo")
}

func (f *archiveKVFetcher) Close() {
	for _, iter := range f.iters {
		iter.Close()
	}
	for _, sst := range f.ssts {
		sst.Close()
	}
}
