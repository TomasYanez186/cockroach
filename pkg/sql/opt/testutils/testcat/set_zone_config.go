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

package testcat

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"gopkg.in/yaml.v2"
)

// SetZoneConfig is a partial implementation of the ALTER TABLE ... CONFIGURE
// ZONE USING statement.
func (tc *Catalog) SetZoneConfig(stmt *tree.SetZoneConfig) *config.ZoneConfig {
	// Update the table name to include catalog and schema if not provided.
	tabName := stmt.TableOrIndex.Table
	tc.qualifyTableName(&tabName)
	tab := tc.Table(&tabName)

	// Handle special case of primary index.
	if stmt.TableOrIndex.Index == "" {
		tab.Indexes[0].IdxZone = makeZoneConfig(stmt.Options)
		return tab.Indexes[0].IdxZone
	}

	for _, idx := range tab.Indexes {
		if idx.IdxName == string(stmt.TableOrIndex.Index) {
			idx.IdxZone = makeZoneConfig(stmt.Options)
			return idx.IdxZone
		}
	}
	panic(fmt.Errorf("\"%q\" is not an index", stmt.TableOrIndex.Index))
}

// makeZoneConfig constructs a ZoneConfig from options provided to the CONFIGURE
// ZONE USING statement.
func makeZoneConfig(options tree.KVOptions) *config.ZoneConfig {
	zone := &config.ZoneConfig{}
	for i := range options {
		switch options[i].Key {
		case "constraints":
			constraintsList := &config.ConstraintsList{}
			value := options[i].Value.(*tree.StrVal).RawString()
			if err := yaml.UnmarshalStrict([]byte(value), constraintsList); err != nil {
				panic(err)
			}
			zone.Constraints = constraintsList.Constraints

		case "lease_preferences":
			value := options[i].Value.(*tree.StrVal).RawString()
			if err := yaml.UnmarshalStrict([]byte(value), &zone.LeasePreferences); err != nil {
				panic(err)
			}
		}
	}
	return zone
}
