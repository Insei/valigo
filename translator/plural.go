package translator

import (
	"strings"
)

// The plural categories of CLDR used by the translations: a key with plural forms is translated
// by the keys key#one, key#few, key#many and key#other, see pluralCategory.
const (
	pluralOne   = "one"
	pluralFew   = "few"
	pluralMany  = "many"
	pluralOther = "other"
)

// isPluralForms reports whether the map of a translation is the plural forms of a message:
// all its keys are the plural categories of CLDR.
func isPluralForms(m map[string]any) bool {
	if len(m) == 0 {
		return false
	}
	for k := range m {
		switch k {
		case "zero", pluralOne, "two", pluralFew, pluralMany, pluralOther:
		default:
			return false
		}
	}
	return true
}

// pluralCategory returns the plural category of CLDR of the integer n in the language lang.
func pluralCategory(lang string, n int64) string {
	if n < 0 {
		n = -n
	}
	base, _, _ := strings.Cut(lang, "-")
	switch base {
	case "ru", "uk", "be":
		mod10, mod100 := n%10, n%100
		switch {
		case mod10 == 1 && mod100 != 11:
			return pluralOne
		case mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14):
			return pluralFew
		default:
			return pluralMany
		}
	case "fr":
		if n <= 1 {
			return pluralOne
		}
		return pluralOther
	default:
		if n == 1 {
			return pluralOne
		}
		return pluralOther
	}
}

// pluralCount returns the first integer of args, the plural form of a message is selected by it.
func pluralCount(args []any) (int64, bool) {
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			return int64(v), true
		case int8:
			return int64(v), true
		case int16:
			return int64(v), true
		case int32:
			return int64(v), true
		case int64:
			return v, true
		case uint:
			return int64(v), true
		case uint8:
			return int64(v), true
		case uint16:
			return int64(v), true
		case uint32:
			return int64(v), true
		case uint64:
			return int64(v), true
		}
	}
	return 0, false
}
