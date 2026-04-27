// Package coerce normalises environment variable values to a target type
// (string, int, bool, float) and reports any keys that could not be converted.
package coerce

import (
	"fmt"
	"strconv"
	"strings"
)

// Type represents the target coercion type.
type Type string

const (
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
	TypeString Type = "string"
)

// Result holds the outcome of a coercion run.
type Result struct {
	Coerced  map[string]string
	Failed   map[string]string // key -> original value that could not be coerced
	Changed  int
}

// Run attempts to coerce each value in vars to the given Type.
// Keys listed in keys are coerced; if keys is empty all keys are processed.
func Run(vars map[string]string, target Type, keys []string) Result {
	filter := toSet(keys)
	res := Result{
		Coerced: make(map[string]string, len(vars)),
		Failed:  make(map[string]string),
	}
	for k, v := range vars {
		if len(filter) > 0 && !filter[k] {
			res.Coerced[k] = v
			continue
		}
		norm, err := coerceValue(v, target)
		if err != nil {
			res.Failed[k] = v
			res.Coerced[k] = v
			continue
		}
		if norm != v {
			res.Changed++
		}
		res.Coerced[k] = norm
	}
	return res
}

// Format returns a human-readable summary of the Result.
func Format(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "coerced %d value(s), %d changed, %d failed\n",
		len(r.Coerced), r.Changed, len(r.Failed))
	for k, v := range r.Failed {
		fmt.Fprintf(&sb, "  FAIL  %s=%q\n", k, v)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func coerceValue(v string, t Type) (string, error) {
	switch t {
	case TypeBool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", err
		}
		return strconv.FormatBool(b), nil
	case TypeInt:
		i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(i, 10), nil
	case TypeFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	default:
		return v, nil
	}
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
