package doer

import (
	"example.com/a/b/external"
)

type (
	Doer interface {
		Do(a, b int, c float64) (d, e int, err error)
	}

	EmbeddedDoer interface {
		Doer
	}

	ExternalEmbeddedDoer interface {
		external.Doer
	}

	doer interface { //nolint:unused
		do(a, b int, c float64) (d, e int, err error)
	}
)
