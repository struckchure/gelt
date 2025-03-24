package profile

type Page struct{}

type Props struct {
	Username string `param:"username"`
}

func (Page) Props() any {
	return Props{}
}
