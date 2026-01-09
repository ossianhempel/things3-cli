package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/things"
)

func authTokenFromEnv() string {
	return strings.TrimSpace(os.Getenv("THINGS_AUTH_TOKEN"))
}

func resolveAuthToken(app *App, explicit string) (string, error) {
	token := strings.TrimSpace(explicit)
	source := "--auth-token"
	if token == "" {
		token = authTokenFromEnv()
		source = "THINGS_AUTH_TOKEN"
	}
	if token == "" {
		return "", things.ErrMissingAuthToken
	}
	if app != nil && app.Debug {
		fmt.Fprintf(app.Err, "Auth token source: %s\n", source)
	}
	return token, nil
}
