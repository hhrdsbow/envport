package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// ToExports converts a map of environment variables into a shell export script.
func ToExports(vars map[string]string) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, vars[k])
	}
	return sb.String()
}
