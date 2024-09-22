package num

//var _ ValueConfigurator[int] = &valueConfigurator[int]{}
//
//type valueConfigurator[T numbers] struct {
//	*baseConfigurator[T]
//}
//
//func (i *valueConfigurator[T]) Required() ValueConfigurator[T] {
//	i.baseConfigurator.Required()
//	return i
//}
//
//func (i *valueConfigurator[T]) AnyOf(allowed ...T) ValueConfigurator[T] {
//	i.baseConfigurator.AnyOf(allowed...)
//	return i
//}
//
//func (i *valueConfigurator[T]) AnyOfInterval(begin, end T) ValueConfigurator[T] {
//	i.baseConfigurator.AnyOfInterval(begin, end)
//	return i
//}
//
//func (i *valueConfigurator[T]) Max(val T) ValueConfigurator[T] {
//	i.baseConfigurator.Max(val)
//	return i
//}
//
//func (i *valueConfigurator[T]) Min(val T) ValueConfigurator[T] {
//	i.baseConfigurator.Min(val)
//	return i
//}
//
//// Custom allows for custom validation logic to be applied to the integer value.
//func (i *valueConfigurator[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) ValueConfigurator[T] {
//	customHelper := shared.NewFieldCustomHelper(i.c.Field, i.c.Helper)
//	i.c.CustomAppend(func(ctx context.Context, h shared.Helper, value any) []shared.Error {
//		return f(ctx, customHelper, value.(*T))
//	})
//	return i
//}
//
//// When allows for conditional validation logic to be applied to the integer value.
//func (i *valueConfigurator[T]) When(whenFn func(ctx context.Context, value *T) bool) ValueConfigurator[T] {
//	if whenFn == nil {
//		return i
//	}
//	base := i.c.NewWithWhen(func(ctx context.Context, value any) bool {
//		v, ok := value.(*T)
//		if !ok {
//			return false
//		}
//		return whenFn(ctx, v)
//	})
//	return &valueConfigurator[T]{
//		&baseConfigurator[T]{
//			base,
//		},
//	}
//}
