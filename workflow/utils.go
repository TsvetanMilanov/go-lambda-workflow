package workflow

import (
	"fmt"
	"strings"
)

func getHandlerKey(method, path string) string {
	// TODO: sanitize path.
	return fmt.Sprintf("%s-%s", strings.ToLower(method), path)
}
