package imported

type (
	Type byte

	Interface interface {
		Imported(Type) Type
	}
)
