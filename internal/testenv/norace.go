//go:build !race

// Package testenv describes the environment of the tests.
package testenv

// RaceEnabled reports whether the race detector is on. Allocations are not stable under it:
// sync.Pool drops items randomly, and the instrumentation of race and coverage moves values to the heap.
const RaceEnabled = false
