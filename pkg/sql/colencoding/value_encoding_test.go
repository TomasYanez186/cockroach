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

package colencoding

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/exec/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/exec/types"
	"github.com/cockroachdb/cockroach/pkg/sql/exec/types/conv"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	semtypes "github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
)

func TestDecodeTableValueToCol(t *testing.T) {
	rng, _ := randutil.NewPseudoRand()
	var buf []byte
	var scratch []byte
	nCols := 1000
	datums := make([]tree.Datum, nCols)
	colTyps := make([]*semtypes.T, nCols)
	typs := make([]types.T, nCols)
	for i := 0; i < nCols; i++ {
		ct := sqlbase.RandType(rng)
		et := conv.FromColumnType(ct)
		if et == types.Unhandled {
			i--
			continue
		}
		datum := sqlbase.RandDatum(rng, ct, false /* nullOk */)
		colTyps[i] = ct
		typs[i] = et
		datums[i] = datum
		var err error
		fmt.Println(datum)
		buf, err = sqlbase.EncodeTableValue(buf, sqlbase.ColumnID(encoding.NoColumnID), datum, scratch)
		if err != nil {
			t.Fatal(err)
		}
	}
	batch := coldata.NewMemBatchWithSize(typs, 1)
	for i := 0; i < nCols; i++ {
		typeOffset, dataOffset, _, typ, err := encoding.DecodeValueTag(buf)
		fmt.Println(typ)
		if err != nil {
			t.Fatal(err)
		}
		buf, err = DecodeTableValueToCol(batch.ColVec(i), 0 /* rowIdx */, typ,
			dataOffset, colTyps[i], buf[typeOffset:])
		if err != nil {
			t.Fatal(err)
		}

		// TODO(jordan): should actually compare the outputs as well, but this is a
		// decent enough smoke test.
	}

	if len(buf) != 0 {
		t.Fatalf("leftover bytes %s", buf)
	}
}
