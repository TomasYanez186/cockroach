// Copyright 2019 The Cockroach Authors.
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

package row

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
)

// fkExistenceCheckForInsert is an auxiliary object that facilitates the existence
// checks on the referenced table when inserting new rows in
// a referencing table.
type fkExistenceCheckForInsert struct {
	// fks maps mutated index id to slice of fkExistenceCheckBaseHelper, the outgoing
	// foreign key existence checkers for each mutated index.
	//
	// In an fkInsertHelper, these slices will have at most one entry,
	// since there can be (currently) at most one outgoing foreign key
	// per mutated index. We use this data structure instead of a
	// one-to-one map for consistency with the other helpers.
	//
	// TODO(knz): this limitation in CockroachDB is arbitrary and
	// incompatible with PostgreSQL. pg supports potentially multiple
	// referencing FK constraints for a single column tuple.
	fks map[sqlbase.IndexID][]fkExistenceCheckBaseHelper

	// checker is the object that actually carries out the lookups in
	// KV.
	checker *fkExistenceBatchChecker
}

// makeFkExistenceCheckHelperForInsert instantiates an insert helper.
func makeFkExistenceCheckHelperForInsert(
	txn *client.Txn,
	table *sqlbase.ImmutableTableDescriptor,
	otherTables FkTableMetadata,
	colMap map[sqlbase.ColumnID]int,
	alloc *sqlbase.DatumAlloc,
) (fkExistenceCheckForInsert, error) {
	h := fkExistenceCheckForInsert{
		checker: &fkExistenceBatchChecker{
			txn: txn,
		},
	}

	// We need an existence check helper for every referenced
	// table. Today, referenced tables are determined by
	// the index descriptors.
	// TODO(knz): make foreign key constraints independent
	// of index definitions.
	for _, idx := range table.AllNonDropIndexes() {
		if idx.ForeignKey.IsSet() {
			fk, err := makeFkExistenceCheckBaseHelper(txn, otherTables, idx, idx.ForeignKey, colMap, alloc, CheckInserts)
			if err == errSkipUnusedFK {
				continue
			}
			if err != nil {
				return h, err
			}
			if h.fks == nil {
				h.fks = make(map[sqlbase.IndexID][]fkExistenceCheckBaseHelper)
			}
			h.fks[idx.ID] = append(h.fks[idx.ID], fk)
		}
	}

	return h, nil
}

// addAllIdxChecks queues a FK existence check for every referenced table.
func (h fkExistenceCheckForInsert) addAllIdxChecks(
	ctx context.Context, row tree.Datums, traceKV bool,
) error {
	for idx := range h.fks {
		if err := queueFkExistenceChecksForRow(ctx, h.checker, h.fks, idx, row, traceKV); err != nil {
			return err
		}
	}
	return nil
}
