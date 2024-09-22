package num

//var _ PtrConfigurator[int] = &ptrConfigurator[int]{}
//
//type ptrConfigurator[T numbers] struct {
//	*baseConfigurator[T]
//}
//
//func (i *ptrConfigurator[T]) Required() PtrConfigurator[T] {
//	i.baseConfigurator.Required()
//	return i
//}
//
//func (i *ptrConfigurator[T]) AnyOf(allowed ...T) PtrConfigurator[T] {
//	i.baseConfigurator.AnyOf(allowed...)
//	return i
//}
//
//func (i *ptrConfigurator[T]) AnyOfInterval(begin, end T) PtrConfigurator[T] {
//	i.baseConfigurator.AnyOfInterval(begin, end)
//	return i
//}
//
//func (i *ptrConfigurator[T]) Max(val T) PtrConfigurator[T] {
//	i.baseConfigurator.Max(val)
//	return i
//}
//
//func (i *ptrConfigurator[T]) Min(val T) PtrConfigurator[T] {
//	i.baseConfigurator.Min(val)
//	return i
//}
//
//// Custom allows for custom validation logic to be applied to the integer value.
//func (i *ptrConfigurator[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) PtrConfigurator[T] {
//	customHelper := shared.NewFieldCustomHelper(i.c.Field, i.c.Helper)
//	i.c.CustomAppend(func(ctx context.Context, h shared.Helper, value any) []shared.Error {
//		return f(ctx, customHelper, value)
//	})
//	return i
//}
//
//// When allows for conditional validation logic to be applied to the integer value.
//func (i *ptrConfigurator[T]) When(whenFn func(ctx context.Context, value any) bool) PtrConfigurator[T] {
//	if whenFn == nil {
//		return i
//	}
//	base := i.c.NewWithWhen(func(ctx context.Context, value any) bool {
//		v, ok := value.(**T)
//		if !ok {
//			return false
//		}
//		return whenFn(ctx, v)
//	})
//	return &ptrConfigurator[T]{
//		&baseConfigurator[T]{
//			base,
//		},
//	}
//}
