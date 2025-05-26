package gelt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/samber/lo"
)

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

func GetComponentName(filename string) string {
	// Get the base name of the file
	base := filepath.Base(filename)

	// Remove the extension
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Convert to CamelCase
	return lo.PascalCase(name)
}

func isNumeric(s string) bool {
	return regexp.MustCompile(`^\d+(\.\d+)?$`).MatchString(s)
}

func isJSExpression(s string) bool {
	// crude: checks for common JS syntax (could be improved)
	return strings.ContainsAny(s, "+-*/()") || strings.HasPrefix(s, "function") || strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[")
}

func MapToJSObject(m map[string]any) string {
	var parts []string

	// Sort keys for deterministic output
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]

		// If value looks like a number or JS expr, output as-is
		if isNumeric(v.(string)) || isJSExpression(v.(string)) {
			parts = append(parts, fmt.Sprintf(`%s: %s`, k, v))
		} else {
			parts = append(parts, fmt.Sprintf(`%s: %q`, k, v)) // Quote value
		}
	}

	return fmt.Sprintf("{ %s }", strings.Join(parts, ", "))
}
