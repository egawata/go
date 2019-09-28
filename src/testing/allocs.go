// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"runtime"
)

// AllocsPerRun returns the average number of allocations during calls to f.
// Although the return value has type float64, it will always be an integral value.
//
// To compute the number of allocations, the function will first be run once as
// a warm-up. The average number of allocations over the specified number of
// runs will then be measured and returned.
//
// AllocsPerRun sets GOMAXPROCS to 1 during its measurement and will restore
// it before returning.
func AllocsPerRun(runs int, f func()) (avg float64) {
	// GOMAXPROCS の戻り値は、これまでの設定値なので、
	// 以下のようにすれば一時的に GOMAXPROCS を変更できる
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

	// Warm up the function
	f()

	// Measure the starting statistics
	// ReadMemStats は現在の stat を返す
	// これを実行後の stat と比較する、という方法
	// Mallocs は object に対して memory allocation が行われた回数
	// (容量ではないらしい)
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	mallocs := 0 - memstats.Mallocs

	// Run the function the specified number of times
	for i := 0; i < runs; i++ {
		f()
	}

	// Read the final statistics
	runtime.ReadMemStats(&memstats)
	mallocs += memstats.Mallocs

	// Average the mallocs over the runs (not counting the warm-up).
	// We are forced to return a float64 because the API is silly, but do
	// the division as integers so we can ask if AllocsPerRun()==1
	// instead of AllocsPerRun()<2.
	return float64(mallocs / uint64(runs))
}
