// Copyright 2014 The Cockroach Authors.
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

package engine

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage/engine/enginepb"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

// MergeInternalTimeSeriesData exports the engine's C++ merge logic for
// InternalTimeSeriesData to higher level packages. This is intended primarily
// for consumption by high level testing of time series functionality.
// If mergeIntoNil is true, then the initial state of the merge is taken to be
// 'nil' and the first operand is merged into nil. If false, the first operand
// is taken to be the initial state of the merge.
// If usePartialMerge is true, the operands are merged together using a partial
// merge operation first, and are then merged in to the initial state. This
// can combine with mergeIntoNil: the initial state is either 'nil' or the first
// operand.
func MergeInternalTimeSeriesData(
	mergeIntoNil, usePartialMerge bool, sources ...roachpb.InternalTimeSeriesData,
) (roachpb.InternalTimeSeriesData, error) {
	// Wrap each proto in an inlined MVCC value, and marshal each wrapped value
	// to bytes. This is the format required by the engine.
	srcBytes := make([][]byte, 0, len(sources))
	var val roachpb.Value
	for _, src := range sources {
		if err := val.SetProto(&src); err != nil {
			return roachpb.InternalTimeSeriesData{}, err
		}
		bytes, err := protoutil.Marshal(&enginepb.MVCCMetadata{
			RawBytes: val.RawBytes,
		})
		if err != nil {
			return roachpb.InternalTimeSeriesData{}, err
		}
		srcBytes = append(srcBytes, bytes)
	}

	// Merge every element into a nil byte slice, one at a time.
	var (
		mergedBytes []byte
		err         error
	)
	if !mergeIntoNil {
		mergedBytes = srcBytes[0]
		srcBytes = srcBytes[1:]
	}
	if usePartialMerge {
		partialBytes := srcBytes[0]
		srcBytes = srcBytes[1:]
		for _, bytes := range srcBytes {
			partialBytes, err = goPartialMerge(partialBytes, bytes)
			if err != nil {
				return roachpb.InternalTimeSeriesData{}, err
			}
		}
		srcBytes = [][]byte{partialBytes}
	}
	for _, bytes := range srcBytes {
		mergedBytes, err = goMerge(mergedBytes, bytes)
		if err != nil {
			return roachpb.InternalTimeSeriesData{}, err
		}
	}

	// Unmarshal merged bytes and extract the time series value within.
	var meta enginepb.MVCCMetadata
	if err := protoutil.Unmarshal(mergedBytes, &meta); err != nil {
		return roachpb.InternalTimeSeriesData{}, err
	}
	mergedTS, err := MakeValue(meta).GetTimeseries()
	if err != nil {
		return roachpb.InternalTimeSeriesData{}, err
	}
	return mergedTS, nil
}
