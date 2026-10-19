package goi18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestCount(t *testing.T) {
	for _, arg := range []any{int(3), int8(3), int16(3), int32(3), int64(3), uint(3), uint8(3), uint16(3), uint32(3), uint64(3)} {
		if n, ok := count([]any{"x", arg}); !ok || n != 3 {
			t.Errorf("count of %T: %d, %v", arg, n, ok)
		}
	}
	if _, ok := count([]any{1.5, "x"}); ok {
		t.Error("a float isn't a count")
	}
}

func TestSetForm(t *testing.T) {
	m := &i18n.Message{}
	for _, c := range []string{"zero", "one", "two", "few", "many", "other"} {
		setForm(m, c, c)
	}
	if *m != (i18n.Message{Zero: "zero", One: "one", Two: "two", Few: "few", Many: "many", Other: "other"}) {
		t.Errorf("forms: %+v", m)
	}
}
