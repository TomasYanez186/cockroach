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

package quotapool

import (
	"context"
	"sort"

	"github.com/cockroachdb/cockroach/pkg/util/ctxgroup"
)

// An example use case for AcquireFunc is a pool of workers attempting to
// acquire resources to run a heterogenous set of jobs. Imagine for example we
// have a set of workers and a list of jobs which need to be run. The function
// might be used to choose the largest job which can be run by the existing
// quantity of quota.
func ExampleIntPool_AcquireFunc() {
	const quota = 7
	const workers = 3
	qp := NewIntPool("work units", quota)
	type job struct {
		name string
		cost int64
	}
	jobs := []*job{
		{name: "foo", cost: 3},
		{name: "bar", cost: 2},
		{name: "baz", cost: 4},
		{name: "qux", cost: 6},
		{name: "quux", cost: 3},
		{name: "quuz", cost: 3},
	}
	// sortJobs sorts the jobs in highest-to-lowest order with nil last.
	sortJobs := func() {
		sort.Slice(jobs, func(i, j int) bool {
			ij, jj := jobs[i], jobs[j]
			if ij != nil && jj != nil {
				return ij.cost > jj.cost
			}
			return ij != nil
		})
	}
	// getJob finds the largest job which can be run with the current quota.
	getJob := func(
		ctx context.Context, qp *IntPool,
	) (j *job, alloc *IntAlloc, err error) {
		alloc, err = qp.AcquireFunc(ctx, func(
			ctx context.Context, v int64,
		) (fulfilled bool, took int64) {
			sortJobs()
			// There are no more jobs, take 0 and return.
			if jobs[0] == nil {
				return true, 0
			}
			// Find the largest jobs which can be run.
			for i := range jobs {
				if jobs[i] == nil {
					break
				}
				if jobs[i].cost <= v {
					j, jobs[i] = jobs[i], nil
					return true, j.cost
				}
			}
			return false, 0
		})
		return j, alloc, err
	}
	runWorker := func(workerNum int) func(ctx context.Context) error {
		return func(ctx context.Context) error {
			for {
				j, alloc, err := getJob(ctx, qp)
				if err != nil {
					return err
				}
				if j == nil {
					return nil
				}
				alloc.Release()
			}
		}
	}
	g := ctxgroup.WithContext(context.Background())
	for i := 0; i < workers; i++ {
		g.GoCtx(runWorker(i))
	}
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
