package gelt

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
)

type Action[T, R any] func(c any) (R, error)

type Load func(c echo.Context) (any, error)

type Handler func() error

type IServer interface {
	Start(...int) error
}

type Server struct {
	srv      *echo.Echo
	router   Router
	registry Registry
}

func (s *Server) Start(port ...int) error {
	if len(port) == 0 {
		port = []int{8080}
	}

	err := s.router.Register(s.srv)
	if err != nil {
		return err
	}

	go s.Watcher()

	return s.srv.Start(fmt.Sprintf(":%d", port[0]))
}

func (s *Server) Watcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Name != "routes_gen.go" {
					err := s.registry.Generate()
					if err != nil {
						log.Fatalln(err)
					}
				}

				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}

				if event.Has(fsnotify.Create) {
					fi, err := os.Stat(event.Name)
					if err == nil && fi.IsDir() {
						log.Println("new directory detected, adding to watcher:", event.Name)
						watcher.Add(event.Name) // Watch new folder
					}

					log.Println("created file:", event.Name)
				}

				if event.Has(fsnotify.Remove) {
					log.Println("removed file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	log.Println("Watching for changes...")

	// Block main goroutine forever.
	<-make(chan struct{})
}

func NewServer(pageRegistry map[string]any) *Server {
	srv := echo.New()

	srv.Use(middleware.Static(JoinURL(lo.Must(os.Getwd()), "public")))

	router := NewRouter(pageRegistry)
	registry := NewRegistry(
		".",
		JoinURL(lo.Must(os.Getwd()),
			"routes_gen.go"),
		lo.Must(GetGoModuleName(".")),
	)

	return &Server{
		srv:      srv,
		router:   *router,
		registry: *registry,
	}
}
