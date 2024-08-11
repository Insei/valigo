package shared

import (
	"context"
)

type FieldValidationFn func(ctx context.Context, h *Helper, v any) []Error
