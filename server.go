package gelt

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Action[T, R any] func(c any) (R, error)

type Load func(c echo.Context) (any, error)

type Handler func() error

type IServer interface {
	Start(...int) error
}

type Server struct {
	srv    *echo.Echo
	router Router
}

func (s *Server) Start(port ...int) error {
	if len(port) == 0 {
		port = []int{8080}
	}

	err := s.router.Register(s.srv)
	if err != nil {
		return err
	}

	return s.srv.Start(fmt.Sprintf(":%d", port[0]))
}

func NewServer(pageRegistry map[string]any) *Server {
	srv := echo.New()

	router := NewRouter(pageRegistry)

	return &Server{
		srv:    srv,
		router: *router,
	}
}
