// This file is part of myhttp

// Licensed under the BSD-3-Clause license:
// https://opensource.org/licenses/BSD-3-Clause
// Copyright (c) 2020, Pablo S. Blum de Aguiar

package main

import (
	"runtime"
)

// requestsPerCPU holds a value
const requestsPerCPU = 8

type limiter interface {
	getNumCPU() int
}

type myLimiter struct{}

var myHTTPLimiter limiter = &myLimiter{}

// getMaxParallel wraps getMaxParallelFromLimiter
func getMaxParallel() uint64 {
	return getMaxParallelFromLimiter(myHTTPLimiter)
}

// getMaxParallel gets an approximation of the currently possible maximum
// value of parallel
func getMaxParallelFromLimiter(l limiter) uint64 {
	return uint64(l.getNumCPU()) * requestsPerCPU
}

// getNumCPU wraps runtime.NumCPU
func (l *myLimiter) getNumCPU() int {
	return runtime.NumCPU()
}
