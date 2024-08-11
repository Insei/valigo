package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

type Helper interface {
	ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) Error
}
