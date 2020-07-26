package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/kelseyhightower/envconfig"
)

// Config of application
type Config struct {
	Server Server
	Gitea  Gitea
	Github Github
}

// Server setting
type Server struct {
	Secret string `default:"secret" envconfig:"GATES_SERVER_SECRET"`
	Addr   string `default:"http://localhost:8080" envconfig:"GATES_SERVER_ADDR"`
	Base   string `envconfig:"GATES_SERVER_BASE"`
}

// Gitea connection setting
type Gitea struct {
	Server       string   `envconfig:"GATES_GITEA_SERVER"`
	ClientID     string   `envconfig:"GATES_GITEA_CLIENT_ID"`
	ClientSecret string   `envconfig:"GATES_GITEA_CLIENT_SECRET"`
	SkipVerity   bool     `envconfig:"GATES_GITEA_SKIP_VERITY"`
	Scope        []string `default:"repo,repo:status,user:email,read:org" envconfig:"GATES_GITEA_SCOPE"`
}

// Github connection setting
type Github struct {
	Server       string   `default:"https://github.com" envconfig:"GATES_GITHUB_SERVER"`
	APIServer    string   `envconfig:"GATES_GITHUB_API_SERVER"`
	ClientID     string   `envconfig:"GATES_GITHUB_CLIENT_ID"`
	ClientSecret string   `envconfig:"GATES_GITHUB_CLIENT_SECRET"`
	SkipVerity   bool     `envconfig:"GATES_GITHUB_SKIP_VERIFY"`
	Scope        []string `default:"repo,repo:status,user:email,read:org" envconfig:"GATES_GITHUB_SCOPE"`
}

func Environ() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	return cfg, err
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
