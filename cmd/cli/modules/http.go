package modules

import (
	"net/http"

	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// GetHTTPClient depends on cli Context.
func GetHTTPClient(c *cli.Context) *http.Client {
	if c.String("token") != "" {
		return oauth2.NewClient(c.Context, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: c.String("token"),
		}))
	}
	return http.DefaultClient
}
