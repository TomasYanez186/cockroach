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

func registerTests(r *registry) {
	// Helpful shell pipeline to generate the list below:
	//
	// grep -h -E 'func register[^(]+\(.*registry\) {' *.go | grep -E -o 'register[^(]+' | grep -E -v '^register(Tests|Benchmarks)$' | grep -v '^\w*Bench$' | sort -f | awk '{printf "\t%s(r)\n", $0}'

	registerAcceptance(r)
	registerAllocator(r)
	registerBackup(r)
	registerCancel(r)
	registerCDC(r)
	registerClearRange(r)
	registerClock(r)
	registerCopy(r)
	registerDecommission(r)
	registerDiskFull(r)
	registerDiskStalledDetection(r)
	registerDrop(r)
	registerElectionAfterRestart(r)
	registerEncryption(r)
	registerFlowable(r)
	registerFollowerReads(r)
	registerGossip(r)
	registerHibernate(r)
	registerHotSpotSplits(r)
	registerImportTPCC(r)
	registerImportTPCH(r)
	registerIndexes(r)
	registerInterleaved(r)
	registerJepsen(r)
	registerKV(r)
	registerKVContention(r)
	registerKVQuiescenceDead(r)
	registerKVGracefulDraining(r)
	registerKVScalability(r)
	registerKVSplits(r)
	registerLargeRange(r)
	registerLedger(r)
	registerNetwork(r)
	registerPsycopg(r)
	registerQueue(r)
	registerRebalanceLoad(r)
	registerReplicaGC(r)
	registerRestore(r)
	registerRoachmart(r)
	registerScaleData(r)
	registerSchemaChangeBulkIngest(r)
	registerSchemaChangeKV(r)
	registerSchemaChangeIndexTPCC100(r)
	registerSchemaChangeIndexTPCC1000(r)
	registerSchemaChangeInvertedIndex(r)
	registerScrubAllChecksTPCC(r)
	registerScrubIndexOnlyTPCC(r)
	registerSyncTest(r)
	registerSysbench(r)
	registerTPCC(r)
	registerTypeORM(r)
	registerLoadSplits(r)
	registerUpgrade(r)
	registerVersion(r)
	registerYCSB(r)
	registerTPCHBench(r)
}

func registerBenchmarks(r *registry) {
	// Helpful shell pipeline to generate the list below:
	//
	// grep -h -E 'func register[^(]+\(.*registry\) {' *.go | grep -E -o 'register[^(]+' | grep -v '^registerTests$' | grep '^\w*Bench$' | sort | awk '{printf "\t%s(r)\n", $0}'

	registerIndexesBench(r)
	registerTPCCBench(r)
	registerTPCHBench(r)
}
