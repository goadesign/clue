package doer

type (
	Doer interface {
		Do(a, b int, c float64) (d, e int, err error)
	}

	doer interface { //nolint:unused
		do(a, b int, c float64) (d, e int, err error)
	}
)
