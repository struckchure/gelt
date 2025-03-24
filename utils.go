package gelt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"strings"
)

func GetFunction[T any](target any, methodName string) (T, bool, error) {
	var zero T // Declare zero value of T to return in case of failure

	// Get the method by name using reflection
	method := reflect.ValueOf(target).MethodByName(methodName)
	if !method.IsValid() {
		return zero, false, errors.New("method not found on the struct")
	}

	// Ensure that T is a function
	fn, ok := method.Interface().(T)
	if !ok {
		return zero, false, errors.New("method signature does not match expected function type")
	}

	return fn, true, nil
}

// GetPackageName extracts the package name from a Go source file
func GetPackageName(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, file, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	return strings.TrimSpace(node.Name.Name), nil
}

func Prettify(v any) (string, error) {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetIndent("", "  ") // Pretty-print with two spaces
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("failed to encode JSON: %w", err)
	}

	return out.String(), nil
}

// ExtractQueryTags extracts all struct field values that have the "query" tag
func ExtractQueryTags(data any) map[string]any {
	result := make(map[string]any)

	// Get the type and value of the struct
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	// If it's a pointer, get the underlying element
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Ensure we're dealing with a struct
	if v.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct")
		return result
	}

	// Iterate over the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get the 'query' tag
		queryTag := field.Tag.Get("query")
		if queryTag != "" {
			result[queryTag] = value.Interface()
		}
	}

	return result
}

func JoinURL(base string, parts ...string) string {
	// Ensure base has no trailing slash
	base = strings.TrimSuffix(base, "/")

	// Join all parts using path.Join (which removes duplicate slashes)
	fullPath := path.Join(append([]string{base}, parts...)...)

	// Ensure the result starts with a slash if the base did
	if strings.HasPrefix(base, "/") {
		fullPath = "/" + strings.TrimPrefix(fullPath, "/")
	}

	return fullPath
}
