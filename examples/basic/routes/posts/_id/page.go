package post_details

import "github.com/labstack/echo/v4"

type Page struct{}

type Data struct {
	Id string
}

func (Page) Load(c echo.Context) (any, error) {
	return &Data{Id: "hello"}, nil
}
