package shared

import (
	"context"
	"testing"
)

// FuzzFieldAccess checks the reading of the fields by the offsets of fmap against the fields themselves:
// nested structs stored by pointer, nil pointers, slices of pointers with nil elements.
func FuzzFieldAccess(f *testing.F) {
	f.Add("text", uint8(3), uint8(0xff))
	f.Add("", uint8(0), uint8(0))
	f.Add("é", uint8(7), uint8(5))
	text, textPtr := sampleField(f, "Text"), sampleField(f, "TextPtr")
	innerText, innerPtr := sampleField(f, "Inner.Text"), sampleField(f, "Inner.TextPtr")
	f.Fuzz(func(t *testing.T, s string, n, flags uint8) {
		obj := &sample{Text: s, Num: int(n)}
		if flags&1 != 0 {
			obj.TextPtr = &s
		}
		if flags&2 != 0 {
			obj.Inner = &inner{Text: s + "i"}
			if flags&4 != 0 {
				obj.Inner.TextPtr = &obj.Text
			}
		}
		for i := range int(n % 8) {
			obj.Texts = append(obj.Texts, s)
			if (int(flags)>>(i%8))&1 != 0 {
				obj.TextPtrs = append(obj.TextPtrs, &obj.Text)
			} else {
				obj.TextPtrs = append(obj.TextPtrs, nil)
			}
		}

		if p := Getter[string](text)(obj); p != &obj.Text {
			t.Errorf("Text: %p, want %p", p, &obj.Text)
		}
		if p := Getter[string](textPtr)(obj); p != obj.TextPtr {
			t.Errorf("TextPtr: %p, want %p", p, obj.TextPtr)
		}
		var wantInner, wantInnerPtr *string
		if obj.Inner != nil {
			wantInner, wantInnerPtr = &obj.Inner.Text, obj.Inner.TextPtr
		}
		if p := Getter[string](innerText)(obj); p != wantInner {
			t.Errorf("Inner.Text: %p, want %p", p, wantInner)
		}
		if p := Getter[string](innerPtr)(obj); p != wantInnerPtr {
			t.Errorf("Inner.TextPtr: %p, want %p", p, wantInnerPtr)
		}
		if got := InStruct(innerText)(obj); got != (obj.Inner != nil) {
			t.Errorf("InStruct: %v", got)
		}

		var r rules
		var got []*string
		NewSliceFieldConfigurator[string](sampleField(t, "TextPtrs"), pathHelper{}, sampleFields, r.add, nil).
			Custom(func(_ context.Context, _ *FieldCustomHelper, e Elems[string]) []Error {
				for _, p := range e.All() {
					got = append(got, p)
				}
				return nil
			})
		r.validate(obj)
		if len(got) != len(obj.TextPtrs) {
			t.Fatalf("elements: %d, want %d", len(got), len(obj.TextPtrs))
		}
		for i, p := range got {
			if p != obj.TextPtrs[i] {
				t.Errorf("element %d: %p, want %p", i, p, obj.TextPtrs[i])
			}
		}
	})
}
