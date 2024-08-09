package helper

import "github.com/insei/valigo/translator"

type Helper struct {
	translator.Translator
}

func NewHelper() *Helper {
	return &Helper{translator.New()}
}
