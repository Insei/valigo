//go:build race

package uuid_test

// raceEnabled reports whether the race detector is on, allocations are not stable under it.
const raceEnabled = true
