package user_by_username

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

type Page struct{}

type names []string

func (n *names) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*n, ", "))
}

type payload struct {
	Names names
}

func (Page) Load(c echo.Context) (any, error) {
	return payload{Names: names{"Bob", "Alice"}}, nil
}
