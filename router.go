package gelt

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type Router struct {
	routes       []Route
	pageRegistry map[string]any
}

type Route struct {
	Pattern  *regexp.Regexp
	Params   []string
	Path     string
	Script   string
	Template string
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
		}

		return nil
	})
}

func (r *Router) matchRoute(path string) (bool, Route, map[string]string) {
	for _, route := range r.routes {
		if matches := route.Pattern.FindStringSubmatch(path); matches != nil {
			params := make(map[string]string)
			for i, name := range route.Params {
				params[name] = matches[i+1]
			}
			return true, route, params
		}
	}
	return false, Route{}, nil
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
	tmpl, err := template.New(filepath.Base(route.Template)).ParseFiles(route.Template)
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

func (r *Router) HandleGet(c echo.Context) error {
	path := r.formatPath(c.Request().URL.Path)

	found, route, _ := r.matchRoute(path)
	if !found {
		return c.HTML(404, "<p>Not found</p>")
	}

	content, err := r.executeRoute(c, route)
	if err != nil {
		return c.HTML(500, fmt.Sprintf("<p>%s</p>", err.Error()))
	}

	return c.HTML(200, *content)
}

func (r *Router) HandlePost(c echo.Context) error {
	return nil
}

func (r *Router) Register(srv *echo.Echo) error {
	err := r.scanPages("./routes")
	if err != nil {
		return err
	}

	srv.GET("*", r.HandleGet)
	srv.POST("*", r.HandlePost)

	return nil
}

func NewRouter(pageRegistry map[string]any) *Router {
	return &Router{
		routes:       []Route{},
		pageRegistry: pageRegistry,
	}
}
