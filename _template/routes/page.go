package index

import "github.com/labstack/echo/v4"

type Page struct{}

type payload struct {
	Title  string
	Stacks []string
}

func (Page) Load(c echo.Context) (any, error) {
	return &payload{
		Title:  "Gelt is amazing!",
		Stacks: []string{"Go", "Echo", "HTML"},
	}, nil
}
