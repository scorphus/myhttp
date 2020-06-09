// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"testing"
)

type myTestLimiter struct {
	numCPU int
	*myLimiter
}

func (l *myTestLimiter) getNumCPU() int {
	return l.numCPU
}

func TestGetMaxParallelFromLimiter(t *testing.T) {
	testCases := []struct {
		numCPU              int
		expectedMaxParallel uint64
	}{
		{4, 32},
		{8, 64},
		{24, 192},
	}
	for _, testCase := range testCases {
		var myHTTPTestLimiter limiter = &myTestLimiter{numCPU: testCase.numCPU}
		maxParallel := getMaxParallelFromLimiter(myHTTPTestLimiter)
		if maxParallel != testCase.expectedMaxParallel {
			t.Errorf("Got %d instead of %d", maxParallel, testCase.expectedMaxParallel)
		}
	}
}
