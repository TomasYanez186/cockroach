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

package tree

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

func TestPlaceholderTypesEquals(t *testing.T) {
	testCases := []struct {
		a, b  PlaceholderTypes
		equal bool
	}{
		{ // 0
			PlaceholderTypes{},
			PlaceholderTypes{},
			true,
		},
		{ // 1
			PlaceholderTypes{types.Int, types.Int},
			PlaceholderTypes{types.Int, types.Int},
			true,
		},
		{ // 2
			PlaceholderTypes{types.Int},
			PlaceholderTypes{types.Int, types.Int},
			false,
		},
		{ // 3
			PlaceholderTypes{types.Int, nil},
			PlaceholderTypes{types.Int, types.Int},
			false,
		},
		{ // 4
			PlaceholderTypes{types.Int, types.Int},
			PlaceholderTypes{types.Int, nil},
			false,
		},
		{ // 5
			PlaceholderTypes{types.Int, nil},
			PlaceholderTypes{types.Int, nil},
			true,
		},
		{ // 6
			PlaceholderTypes{types.Int},
			PlaceholderTypes{types.Int, nil},
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res := tc.a.Equals(tc.b)
			if res != tc.equal {
				t.Errorf("%v vs %v: expected %t, got %t", tc.a, tc.b, tc.equal, res)
			}
			res2 := tc.b.Equals(tc.a)
			if res != res2 {
				t.Errorf("%v vs %v: not commutative", tc.a, tc.b)
			}
		})
	}
}
