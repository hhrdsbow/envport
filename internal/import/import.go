// Package import provides functionality to import environment variables
// from external file formats (shell exports, dotenv, JSON) into a snapshot.
package import

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents the input file format.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
)

// Parse reads environment variables from r in the given format and returns
// a map of key/value pairs. Returns an error if the format is unknown or
// the content is malformed.
func Parse(r io.Reader, format Format) (map[string]string, error) {
	switch format {
	case FormatShell:
		return parseShell(r)
	case FormatDotenv:
		return parseDotenv(r)
	case FormatJSON:
		return parseJSON(r)
	default:
		return nil, fmt.Errorf("unknown import format: %q", format)
	}
}

// parseShell handles lines like: export KEY=VALUE or KEY=VALUE
func parseShell(r io.Reader) (map[string]string, error) {
	vars := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		vars[strings.TrimSpace(k)] = strings.Trim(strings.TrimSpace(v), "\"'")
	}
	return vars, scanner.Err()
}

// parseDotenv handles KEY=VALUE lines, stripping quotes and comments.
func parseDotenv(r io.Reader) (map[string]string, error) {
	return parseShell(r) // dotenv format is compatible with our shell parser
}

// parseJSON handles a flat JSON object: {"KEY": "VALUE", ...}
func parseJSON(r io.Reader) (map[string]string, error) {
	vars := make(map[string]string)
	if err := json.NewDecoder(r).Decode(&vars); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return vars, nil
}
