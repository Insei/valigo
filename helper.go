package valigo

import "github.com/insei/valigo/translator"

type Helper struct {
	translator.Translator
}

func newHelper() *Helper {
	return &Helper{translator.New()}
}
