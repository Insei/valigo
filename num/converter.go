package num

type ptrConverter[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64] struct {
}

func numPtrDereference[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](val *T) T {
	return *val
}

func numPtrPtrDereference[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](val **T) T {
	return **val
}
