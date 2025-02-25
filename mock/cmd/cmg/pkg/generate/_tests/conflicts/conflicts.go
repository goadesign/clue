package testing

type (
	Conflicts interface {
		Simple(c *Conflicts) *Conflicts
		AddSimple()
		SetSimple()
		HasMore()
		HasMoreMock()
	}
)
