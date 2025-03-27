package profile

import "github.com/labstack/echo/v4"

type Page struct{}

type Props struct {
	Username string `param:"username"`
}

func (Page) Props() any {
	return Props{}
}

type payload struct {
	Name    string
	Project string
}

func (Page) Load(c echo.Context) (any, error) {
	return payload{Name: "Mohammed", Project: "Gelt"}, nil
}
