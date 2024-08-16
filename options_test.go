package valigo

import (
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/translator"
)

const (
	customRegexpLocaleMsg = "Only numbers and words is allowed"
	customRegexpLocaleKey = "validation:string:" + customRegexpLocaleMsg
)

func TestOptionApply(t *testing.T) {
	tStorage := translator.NewInMemStorage()
	tStorage.AddTranslations("en", map[string]string{
		customRegexpLocaleKey: customRegexpLocaleMsg,
	})
	tr := translator.New(translator.WithStorage(tStorage), translator.WithDefaultLang("en"))
	opt := WithTranslator(tr)
	v := New(opt)
	opt.apply(v)
	if v.helper.t == nil {
		t.Errorf("expected translator to be set")
	}
}

func TestOptionApplyWithNilTranslator(t *testing.T) {
	opt := WithTranslator(nil)
	v := New(opt)
	opt.apply(v)
	if v.helper.t == nil {
		t.Errorf("expected translator to be nil")
	}
}

func TestOptionApplyWithFieldLocationNamingFn(t *testing.T) {
	fn := func(field fmap.Field) string {
		return "test"
	}
	opt := WithFieldLocationNamingFn(fn)
	v := New(opt)
	v.helper.getFieldLocation = nil
	opt.apply(v)
	if v.helper.getFieldLocation == nil {
		t.Errorf("expected getFieldLocation to be set")
	}
}

func TestOptionApplyWithNilFieldLocationNamingFn(t *testing.T) {
	opt := WithFieldLocationNamingFn(nil)
	v := New(opt)
	opt.apply(v)
	if v.helper.getFieldLocation == nil {
		t.Errorf("expected getFieldLocation to be nil")
	}
}

func TestOptionApplyWithMultipleOptions(t *testing.T) {
	tStorage := translator.NewInMemStorage()
	tStorage.AddTranslations("en", map[string]string{
		customRegexpLocaleKey: customRegexpLocaleMsg,
	})
	tr := translator.New(translator.WithStorage(tStorage), translator.WithDefaultLang("en"))
	opt1 := WithTranslator(tr)
	opt2 := WithFieldLocationNamingFn(func(field fmap.Field) string {
		return "test"
	})
	v := New(opt1, opt2)
	opt1.apply(v)
	opt2.apply(v)
	if v.helper.t == nil {
		t.Errorf("expected translator to be set")
	}
	if v.helper.getFieldLocation == nil {
		t.Errorf("expected getFieldLocation to be set")
	}
}
