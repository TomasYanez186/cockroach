// Copyright 2017 The Cockroach Authors.
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

package tree

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

func TestCastToCollatedString(t *testing.T) {
	cases := []struct {
		typ      *types.T
		contents string
	}{
		{types.MakeCollatedString(types.String, "de"), "test"},
		{types.MakeCollatedString(types.String, "en"), "test"},
		{types.MakeCollatedString(types.MakeString(5), "en"), "test"},
		{types.MakeCollatedString(types.MakeString(4), "en"), "test"},
		{types.MakeCollatedString(types.MakeString(3), "en"), "tes"},
	}
	for _, cas := range cases {
		t.Run("", func(t *testing.T) {
			expr := &CastExpr{Expr: NewDString("test"), Type: cas.typ, SyntaxMode: CastShort}
			typedexpr, err := expr.TypeCheck(&SemaContext{}, types.Any)
			if err != nil {
				t.Fatal(err)
			}
			evalCtx := NewTestingEvalContext(cluster.MakeTestingClusterSettings())
			defer evalCtx.Stop(context.Background())
			val, err := typedexpr.Eval(evalCtx)
			if err != nil {
				t.Fatal(err)
			}
			switch v := val.(type) {
			case *DCollatedString:
				if v.Locale != cas.typ.Locale() {
					t.Errorf("expected locale %q but got %q", cas.typ.Locale(), v.Locale)
				}
				if v.Contents != cas.contents {
					t.Errorf("expected contents %q but got %q", cas.contents, v.Contents)
				}
			default:
				t.Errorf("expected type *DCollatedString but got %T", v)
			}
		})
	}
}
