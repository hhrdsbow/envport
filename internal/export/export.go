package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
)

// Write writes environment variables to w in the specified format.
func Write(w io.Writer, env map[string]string, format Format) error {
	keys := sortedKeys(env)
	switch format {
	case FormatShell:
		return writeShell(w, env, keys)
	case FormatDotenv:
		return writeDotenv(w, env, keys)
	case FormatJSON:
		return writeJSON(w, env, keys)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func writeShell(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%q\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeDotenv(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		val := strings.ReplaceAll(env[k], "\n", "\\n")
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, val); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, env map[string]string, keys []string) error {
	fmt.Fprint(w, "{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		val := strings.ReplaceAll(env[k], `"`, `\"`)
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", k, val, comma); err != nil {
			return err
		}
	}
	fmt.Fprint(w, "}\n")
	return nil
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
