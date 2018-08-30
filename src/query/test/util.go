// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/m3db/m3/src/query/block"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// EqualsWithNans helps compare float slices which have NaNs in them
func EqualsWithNans(t *testing.T, expected interface{}, actual interface{}) {
	EqualsWithNansWithDelta(t, expected, actual, 0)
}

// EqualsWithNansWithDelta helps compare float slices which have NaNs in them
// allowing a delta for float comparisons.
func EqualsWithNansWithDelta(t *testing.T, expected interface{}, actual interface{}, delta float64) {
	debugMsg := fmt.Sprintf("expected: %v, actual: %v", expected, actual)
	switch v := expected.(type) {
	case [][]float64:
		actualV, ok := actual.([][]float64)
		require.True(t, ok, "actual should be of type [][]float64, found: %T", actual)
		require.Equal(t, len(v), len(actualV),
			fmt.Sprintf("expected length: %v, actual length: %v\nfor expected: %v, actual: %v",
				len(v), len(actualV), expected, actual))
		for i, vals := range v {
			equalsWithNans(t, vals, actualV[i], delta, debugMsg)
		}

	case []float64:
		actualV, ok := actual.([]float64)
		require.True(t, ok, "actual should be of type []float64, found: %T", actual)
		require.Equal(t, len(v), len(actualV),
			fmt.Sprintf("expected length: %v, actual length: %v\nfor expected: %v, actual: %v",
				len(v), len(actualV), expected, actual))
		equalsWithNans(t, v, actualV, delta, debugMsg)

	case float64:
		actualV, ok := actual.(float64)
		require.True(t, ok, "actual should be of type float64, found: %T", actual)
		equalsWithDelta(t, v, actualV, delta, debugMsg)

	default:
		require.Fail(t, "unknown type: %T", v)
	}
}

func equalsWithNans(t *testing.T, expected []float64, actual []float64, delta float64, debugMsg string) {
	require.Equal(t, len(expected), len(actual))
	for i, v := range expected {
		if math.IsNaN(v) {
			require.True(t, math.IsNaN(actual[i]), debugMsg)
		} else {
			equalsWithDelta(t, v, actual[i], delta, debugMsg)
		}
	}
}

func equalsWithDelta(t *testing.T, expected, actual, delta float64, debugMsg string) {
	if math.IsNaN(expected) {
		require.True(t, math.IsNaN(actual), debugMsg)
	} else {
		if delta == 0 {
			require.Equal(t, expected, actual, debugMsg)
		} else {
			diff := math.Abs(expected - actual)
			require.True(t, delta > diff, debugMsg)
			require.Equal(t, v, actual[i], debugMsg)
		}
	}
}

type match struct {
	indices []int
	metas   block.SeriesMeta
	values  []float64
}

type matches []match

func (m matches) Len() int           { return len(m) }
func (m matches) Less(i, j int) bool { return m[i].metas.Tags.ID() > m[j].metas.Tags.ID() }
func (m matches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

// CompareLists compares series meta / index pairs
func CompareLists(t *testing.T, meta, exMeta []block.SeriesMeta, index, exIndex [][]int) {
	require.Equal(t, len(exIndex), len(exMeta))
	require.Equal(t, len(exMeta), len(meta))
	require.Equal(t, len(exIndex), len(index))

	ex := make(matches, len(meta))
	actual := make(matches, len(meta))
	// build matchers
	for i := range meta {
		ex[i] = match{exIndex[i], exMeta[i], []float64{}}
		actual[i] = match{index[i], meta[i], []float64{}}
	}
	sort.Sort(ex)
	sort.Sort(actual)
	assert.Equal(t, ex, actual)
}

// CompareValues compares series meta / value pairs
func CompareValues(t *testing.T, meta, exMeta []block.SeriesMeta, vals, exVals [][]float64) {
	require.Equal(t, len(exVals), len(exMeta))
	require.Equal(t, len(exMeta), len(meta))
	require.Equal(t, len(exVals), len(vals))

	ex := make(matches, len(meta))
	actual := make(matches, len(meta))
	// build matchers
	for i := range meta {
		ex[i] = match{[]int{}, exMeta[i], exVals[i]}
		actual[i] = match{[]int{}, meta[i], vals[i]}
	}

	sort.Sort(ex)
	sort.Sort(actual)
	for i := range ex {
		assert.Equal(t, ex[i].metas, actual[i].metas)
		EqualsWithNansWithDelta(t, ex[i].values, actual[i].values, 0.00001)
	}
}
