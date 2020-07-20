package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/code-devel-cover/CodeCover/core"
)

// Config of application
type Config struct {
	Server Server
	Gitea  Gitea
	Github Github
}

// Server setting
type Server struct {
	Secret string `default:"secret"`
	Addr   string
	Base   string
}

// Gitea connection setting
type Gitea struct {
	Server       string
	ClientID     string
	ClientSecret string
	SkipVerity   bool
	Scope        []string `default:"repo,repo:status,user:email,read:org"`
}

// Github connection setting
type Github struct {
	Server       string `default:"https://github.com"`
	APIServer    string
	ClientID     string
	ClientSecret string
	SkipVerity   bool
	Scope        []string `default:"repo,repo:status,user:email,read:org"`
}

// Providers of all available SCM
func (c *Config) Providers() []core.SCMProvider {
	providers := make([]core.SCMProvider, 0)
	if c.Gitea.Server != "" {
		providers = append(providers, core.Gitea)
	}
	if c.Github.Server != "" {
		providers = append(providers, core.Github)
	}
	return providers
}

// Port opened for the current server
func (server Server) Port() string {
	u, err := url.Parse(server.Addr)
	if err != nil {
		return ""
	}
	return u.Port()
}

// BaseURL sanitize the Base string to URL format
func (server Server) BaseURL() string {
	base := strings.Trim(server.Base, "/")
	return fmt.Sprintf("/%s", base)
}
