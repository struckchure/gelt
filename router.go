package gelt

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"golang.org/x/net/websocket"
)

type Router struct {
	routes       []Route
	pageRegistry map[string]any
	eventBus     *EventBus
}

type Route struct {
	Pattern  *regexp.Regexp
	Params   []string
	Path     string
	Script   string
	Template string
	Weight   int
}

// ComputeWeight calculates the weight of a route
func ComputeWeight(path string) int {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	n := len(segments)
	weight := 0
	hasParamInFirstSegment := strings.HasPrefix(segments[0], ":")

	for i, segment := range segments {
		if strings.HasPrefix(segment, ":") { // It's a parameter
			weight += (n - i) * 10 // Increase weight for params appearing earlier
		}
	}

	// If first segment is a param, reduce priority
	if hasParamInFirstSegment {
		weight -= 50
	}

	return weight
}

// SortRoutes sorts routes based on their weight
func SortRoutes(routes []Route) {
	for i := range routes {
		routes[i].Weight = ComputeWeight(routes[i].Path)
	}

	sort.SliceStable(routes, func(i, j int) bool {
		// Static routes (weight = 0) always come first
		if routes[i].Weight == 0 && routes[j].Weight > 0 {
			return true
		}
		if routes[j].Weight == 0 && routes[i].Weight > 0 {
			return false
		}

		// Non-static routes sorted from highest to lowest weight
		return routes[i].Weight > routes[j].Weight
	})
}

func (r *Router) scanPages(root string) error {
	paths := make(map[string]bool)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)
		if filename == "page.go" || filename == "page.html" {
			dir := filepath.Dir(path)
			// Get the relative path from root directory
			routePath, err := filepath.Rel(root, dir)
			if err != nil {
				return err
			}
			routePath = strings.Trim(routePath, ".")

			// Skip if we already processed this path
			if paths[routePath] {
				return nil
			}

			paths[routePath] = true

			scriptPath := filepath.Join(dir, "page.go")
			templatePath := filepath.Join(dir, "page.html")

			// Verify both files exist
			if _, err := os.Stat(scriptPath); err != nil {
				scriptPath = ""
			}
			if _, err := os.Stat(templatePath); err != nil {
				templatePath = ""
			}

			route := Route{
				Pattern:  nil,
				Path:     routePath,
				Script:   scriptPath,
				Template: templatePath,
			}
			route.Path = r.getRoutePath(route)

			paramNames := []string{}

			regexPattern := regexp.MustCompile(`:[a-zA-Z0-9_]+`)
			matches := regexPattern.FindAllString(route.Path, -1)

			for _, match := range matches {
				paramNames = append(paramNames, match[1:]) // Remove ':'
			}

			re := `^` + regexPattern.ReplaceAllString(route.Path, `([^/]+)`) + `$`

			route.Pattern = regexp.MustCompile(re)
			route.Params = paramNames

			r.routes = append(r.routes, route)
			SortRoutes(r.routes)
		}

		return nil
	})
}

func (r *Router) executeRoute(ctx echo.Context, route Route) (*string, error) {
	var data any
	var loadFn func(echo.Context) (any, error) = nil

	if route.Script != "" {
		pkgName, err := GetPackageName(route.Script)
		if err != nil {
			return nil, err
		}

		page := r.pageRegistry[pkgName]

		loadFunc, ok, err := GetFunction[any](page, "Load")
		if err == nil {
			if ok {
				loadFn = loadFunc.(func(echo.Context) (any, error))
			} else {
				return nil, err
			}
		}
	}

	if loadFn != nil {
		_data, err := loadFn(ctx)
		if err != nil {
			return nil, err
		}

		data = _data
	}

	// Parse and execute the template with the loaded data
	tmpl, err := template.
		New(filepath.Base(route.Template)).
		Funcs(template.FuncMap{
			"json": ToJSON,
			"sub":  func(a, b int) int { return a - b },
			"add":  func(a, b int) int { return a + b },
			"div":  func(a, b int) int { return a / b },
			"mul":  func(a, b int) int { return a * b },
			"mod":  func(a, b int) int { return a % b },
		}).
		ParseFiles(route.Template)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, data)
	if err != nil {
		return nil, err
	}

	content := builder.String()
	return &content, nil
}

func (r *Router) getRoutePath(route Route) string {
	route.Path = r.formatPath(route.Path)

	segments := strings.Split(strings.Trim(route.Path, "/"), "/")

	for idx, segment := range segments {
		if strings.HasPrefix(segment, "_") {
			segments[idx] = ":" + segment[1:]
			continue
		}
	}

	return r.formatPath(strings.Join(segments, "/"))
}

func (r *Router) AddRoute(path string, route Route) error {
	r.routes = append(r.routes, route)

	return nil
}

func (r *Router) formatPath(path string) string {
	path = strings.Trim(path, "/")

	// add trailing and leading slashes
	if path == "" {
		path = "/"
	}

	path = lo.Ternary(strings.HasPrefix(path, "/"), path, "/"+path)
	path = lo.Ternary(strings.HasSuffix(path, "/"), path, path+"/")

	return path
}

func (r *Router) Register(srv *echo.Echo) error {
	err := r.scanPages("./routes")
	if err != nil {
		return err
	}

	handler := func(c echo.Context, route Route) error {
		content, err := r.executeRoute(c, route)
		if err != nil {
			log.Println(err)
			return c.HTML(500, fmt.Sprintf("<p>%s</p>", err.Error()))
		}

		return c.HTML(200, *content)
	}

	wsHandler := func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			eventChan := r.eventBus.Subscribe("file_changed")

			// Handle sending events
			done := make(chan struct{})
			go func() {
				for {
					select {
					case _, ok := <-eventChan:
						if !ok {
							return // Channel closed, exit goroutine
						}
						err := websocket.Message.Send(ws, "reload")
						if err != nil {
							log.Println("WebSocket Send Error:", err)
							close(done) // Signal to close
							return
						}
						log.Println("Reloading web page...")
					case <-done:
						return
					}
				}
			}()

			// Block until WebSocket closes
			buf := make([]byte, 1)
			_, err := ws.Read(buf)
			close(done) // Ensure goroutine exits
			if err != nil {
				log.Println("WebSocket Read Error:", err)
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	}

	rootHandler := func(c echo.Context) error {
		content, err := r.executeRoute(c, Route{
			Path:     "/",
			Template: "public/index.html",
		})
		if err != nil {
			log.Println(err)
			return c.HTML(500, fmt.Sprintf("<p>%s</p>", err.Error()))
		}

		return c.HTML(200, *content)
	}

	srv.GET("/", rootHandler)
	srv.POST("/", rootHandler)

	srv.GET("/_/ws/", wsHandler)

	for _, route := range r.routes {
		_handler := func(c echo.Context) error { return handler(c, route) }
		srv.Group("_").GET(route.Path, _handler)
		srv.Group("_").POST(route.Path, _handler)
	}

	srv.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == http.StatusNotFound {
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			}
		}
		// Default error handler
		srv.DefaultHTTPErrorHandler(err, c)
	}

	return nil
}

func NewRouter(pageRegistry map[string]any, eventBus *EventBus) *Router {
	return &Router{
		routes:       []Route{},
		pageRegistry: pageRegistry,
		eventBus:     eventBus,
	}
}
