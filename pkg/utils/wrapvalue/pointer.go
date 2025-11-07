package wv

func ValueOrDefault[T any](pointer *T, defaultValue T) T {
	if pointer == nil {
		return defaultValue
	}
	return *pointer
}

func ToPointer[T any](value T) *T {
	return &value
}
