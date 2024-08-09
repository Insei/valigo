package valigo

import (
	"context"
	"fmt"
)

type Helper struct {
	Translator
}

func newHelper() *Helper {
	return &Helper{newNoopTranslator()}
}

type translator struct {
}

func (t *translator) ErrorT(_ context.Context, format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func (t *translator) T(_ context.Context, format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func newNoopTranslator() Translator {
	return &translator{}
}
