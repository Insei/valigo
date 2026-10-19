package timex

import (
	"testing"
	"time"
)

// SetNow replaces the current time of InPast and InFuture for the test.
func SetNow(t testing.TB, fn func() time.Time) {
	now = fn
	t.Cleanup(func() { now = time.Now })
}
