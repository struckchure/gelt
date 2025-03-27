package post_list

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type Page struct{}

type Data struct {
	Page  int
	Size  int
	Posts []*string
}

func (Page) Load(c echo.Context) (any, error) {
	return &Data{
		Page:  1,
		Size:  10,
		Posts: []*string{lo.ToPtr("one")},
	}, nil
}
