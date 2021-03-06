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

package main

import (
	"context"
	gosql "database/sql"
	"fmt"
	"time"
)

// isAlive returns whether the node queried by db is alive
func isAlive(db *gosql.DB) bool {
	_, err := db.Exec("SHOW DATABASES")
	return err == nil
}

// dbUnixEpoch returns the current time in db
func dbUnixEpoch(db *gosql.DB) (float64, error) {
	var epoch float64
	if err := db.QueryRow("SELECT now()::DECIMAL").Scan(&epoch); err != nil {
		return 0, err
	}
	return epoch, nil
}

// offsetInjector is used to inject clock offsets in roachtests
type offsetInjector struct {
	c        *cluster
	deployed bool
}

// deploy installs ntp and downloads / compiles bumptime used to create a clock offset
func (oi *offsetInjector) deploy(ctx context.Context) {
	if err := oi.c.RunE(ctx, oi.c.All(), "test -x ./bumptime"); err != nil {
		oi.c.Install(ctx, oi.c.All(), "ntp")
		oi.c.Install(ctx, oi.c.All(), "gcc")
		oi.c.Run(ctx, oi.c.All(), "sudo", "service", "ntp", "stop")
		oi.c.Run(ctx,
			oi.c.All(),
			"curl",
			"-kO",
			"https://raw.githubusercontent.com/cockroachdb/jepsen/master/cockroachdb/resources/bumptime.c",
		)
		oi.c.Run(ctx, oi.c.All(), "gcc", "bumptime.c", "-o", "bumptime", "&&", "rm bumptime.c")
	}

	oi.deployed = true
}

// offset injects a offset of s into the node with the given nodeID
func (oi *offsetInjector) offset(ctx context.Context, nodeID int, s time.Duration) {
	if !oi.deployed {
		oi.c.t.Fatal("Offset injector must be deployed before injecting a clock offset")
	}

	oi.c.Run(
		ctx,
		oi.c.Node(nodeID),
		fmt.Sprintf("sudo ./bumptime %f", float64(s)/float64(time.Millisecond)),
	)
}

// recover force syncs time on the node with the given nodeID to recover
// from any offsets
func (oi *offsetInjector) recover(ctx context.Context, nodeID int) {
	if !oi.deployed {
		oi.c.t.Fatal("Offset injector must be deployed before recovering from clock offsets")
	}

	syncCmds := [][]string{
		{"sudo", "service", "ntp", "stop"},
		{"sudo", "ntpdate", "-u", "time.google.com"},
		{"sudo", "service", "ntp", "start"},
	}
	for _, cmd := range syncCmds {
		oi.c.Run(
			ctx,
			oi.c.Node(nodeID),
			cmd...,
		)
	}
}

// newOffsetInjector creates a offsetInjector which can be used to inject
// and recover from clock offsets
func newOffsetInjector(c *cluster) *offsetInjector {
	return &offsetInjector{c: c}
}
