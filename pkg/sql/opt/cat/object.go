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

package cat

// Object is implemented by all objects in the catalog.
type Object interface {
	// ID is the unique, stable identifier for this object. See the comment for
	// StableID for more detail.
	ID() StableID

	// Equals returns true if this object is identical to the given Object.
	//
	// Two objects are identical if they have the same identifier and there were
	// no changes to schema or table statistics between the times the two objects
	// were resolved.
	//
	// Used for invalidating cached plans.
	Equals(other Object) bool
}
