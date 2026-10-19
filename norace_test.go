//go:build !race

package valigo

// raceEnabled reports whether the race detector is on, sync.Pool drops items randomly under it.
const raceEnabled = false
