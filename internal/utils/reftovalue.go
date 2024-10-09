package utils

func FromPtr[T any](ref *T) (value T) {
	if ref != nil {
		value = *ref
	}
	return
}
