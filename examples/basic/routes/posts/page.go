package post_list

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type Page struct{}

func (Page) Load(c echo.Context) (any, error) {
	return &Props{
		Page:  1,
		Size:  10,
		Posts: []*string{lo.ToPtr("one")},
	}, nil
}

type Props struct {
	Page  int `query:"page"`
	Size  int `query:"size"`
	Posts []*string
}

func (Page) Props() any {
	return Props{}
}
