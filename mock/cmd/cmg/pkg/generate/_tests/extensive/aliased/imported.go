package imported

type (
	Type byte

	Interface interface {
		Imported(Type) Type
	}

	Generic[T any] interface {
		ImportedGeneric(T) T
	}
)
