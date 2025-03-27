package post_create

import "github.com/labstack/echo/v4"

type Props struct {
	User  string
	Posts []*string
}

type Page struct{}

func (Page) Load(c echo.Context) (any, error) {
	return &Props{
		User:  "new-user",
		Posts: []*string{},
	}, nil
}

type ContactForm struct {
	Name    string
	Email   string
	Message string
}

type ContactFormResponse struct {
	Status bool
}

func (Page) Action(c echo.Context) (*ContactFormResponse, error) {
	return &ContactFormResponse{}, nil
}
