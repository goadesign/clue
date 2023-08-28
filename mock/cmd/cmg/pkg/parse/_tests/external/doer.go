package external

type (
	Doer interface {
		Do(a, b int, c float64) (d, e int, err error)
	}
)
