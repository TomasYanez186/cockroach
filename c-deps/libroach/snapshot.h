// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License included
// in the file licenses/BSL.txt and at www.mariadb.com/bsl11.
//
// Change Date: 2022-10-01
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by the Apache License, Version 2.0,
// included in the file licenses/APL.txt and at
// https://www.apache.org/licenses/LICENSE-2.0

#pragma once

#include <libroach.h>
#include <rocksdb/db.h>
#include "engine.h"

namespace cockroach {

struct DBSnapshot : public DBEngine {
  const rocksdb::Snapshot* snapshot;

  DBSnapshot(DBEngine* db) : DBEngine(db->rep, db->iters), snapshot(db->rep->GetSnapshot()) {}
  virtual ~DBSnapshot();

  virtual DBStatus Put(DBKey key, DBSlice value);
  virtual DBStatus Merge(DBKey key, DBSlice value);
  virtual DBStatus Delete(DBKey key);
  virtual DBStatus SingleDelete(DBKey key);
  virtual DBStatus DeleteRange(DBKey start, DBKey end);
  virtual DBStatus CommitBatch(bool sync);
  virtual DBStatus ApplyBatchRepr(DBSlice repr, bool sync);
  virtual DBSlice BatchRepr();
  virtual DBStatus Get(DBKey key, DBString* value);
  virtual DBIterator* NewIter(DBIterOptions);
  virtual DBStatus GetStats(DBStatsResult* stats);
  virtual DBStatus GetTickersAndHistograms(DBTickersAndHistogramsResult* stats);
  virtual DBString GetCompactionStats();
  virtual DBStatus GetEnvStats(DBEnvStatsResult* stats);
  virtual DBStatus GetEncryptionRegistries(DBEncryptionRegistries* result);
  virtual DBStatus EnvWriteFile(DBSlice path, DBSlice contents);
  virtual DBStatus EnvOpenFile(DBSlice path, rocksdb::WritableFile** file);
  virtual DBStatus EnvReadFile(DBSlice path, DBSlice* contents);
  virtual DBStatus EnvAppendFile(rocksdb::WritableFile* file, DBSlice contents);
  virtual DBStatus EnvSyncFile(rocksdb::WritableFile* file);
  virtual DBStatus EnvCloseFile(rocksdb::WritableFile* file);
  virtual DBStatus EnvDeleteFile(DBSlice path);
  virtual DBStatus EnvDeleteDirAndFiles(DBSlice dir);
  virtual DBStatus EnvLinkFile(DBSlice oldname, DBSlice newname);
};

}  // namespace cockroach
