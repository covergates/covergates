package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/covergates/covergates/core"
	"github.com/kelseyhightower/envconfig"
)

// Config of application
type Config struct {
	Server    Server
	Gitea     Gitea
	Github    Github
	GitLab    GitLab
	Bitbucket Bitbucket
	Database  Database
	CloudRun  CloudRun
}

// Server setting
type Server struct {
	Secret      string `default:"secret" envconfig:"GATES_SERVER_SECRET"`
	Addr        string `default:"http://localhost:8080" envconfig:"GATES_SERVER_ADDR"`
	Base        string `envconfig:"GATES_SERVER_BASE"`
	SkipVerity  bool   `default:"true" envconfig:"GATES_SERVER_SKIP_VERIFY"`
	ServerPort  string `envconfig:"GATES_SERVER_PORT"`
	CloudPort   string `envconfig:"PORT"`
	OAuthClient string `default:"client"`
}

// Database setting
type Database struct {
	AutoMigrate bool   `default:"true" envconfig:"GATES_DB_AUTO_MIGRATE"`
	Driver      string `default:"sqlite3" envconfig:"GATES_DB_DRIVER"`
	Host        string `envconfig:"GATES_DB_HOST"`
	Port        string `envconfig:"GATES_DB_PORT"`
	User        string `envconfig:"GATES_DB_USER"`
	Name        string `default:"core.db" envconfig:"GATES_DB_NAME"`
	Password    string `envconfig:"GATES_DB_PASSWORD"`
}

// CloudRun database setting for google cloud run
type CloudRun struct {
	User     string `envconfig:"GATES_DB_USER"`
	Password string `envconfig:"GATES_DB_PASSWORD"`
	Socket   string `default:"/cloudsql" envconfig:"GATES_DB_SOCKET_DIR"`
	Instance string `envconfig:"INSTANCE_CONNECTION_NAME"`
	Name     string `envconfig:"GATES_DB_NAME"`
}

// Gitea connection setting
type Gitea struct {
	Server       string   `envconfig:"GATES_GITEA_SERVER"`
	ClientID     string   `envconfig:"GATES_GITEA_CLIENT_ID"`
	ClientSecret string   `envconfig:"GATES_GITEA_CLIENT_SECRET"`
	SkipVerity   bool     `envconfig:"GATES_GITEA_SKIP_VERIFY"`
	Scope        []string `default:"repo,repo:status,user:email,read:org" envconfig:"GATES_GITEA_SCOPE"`
}

// Github connection setting
type Github struct {
	Server       string   `default:"https://github.com" envconfig:"GATES_GITHUB_SERVER"`
	APIServer    string   `default:"https://api.github.com" envconfig:"GATES_GITHUB_API_SERVER"`
	ClientID     string   `envconfig:"GATES_GITHUB_CLIENT_ID"`
	ClientSecret string   `envconfig:"GATES_GITHUB_CLIENT_SECRET"`
	SkipVerity   bool     `envconfig:"GATES_GITHUB_SKIP_VERIFY"`
	Scope        []string `default:"repo,repo:status,user:email,read:org" envconfig:"GATES_GITHUB_SCOPE"`
}

// Bitbucket connection settings
type Bitbucket struct {
	Server       string   `default:"https://bitbucket.com" envconfig:"GATES_BITBUCKET_SERVER"`
	ClientID     string   `envconfig:"GATES_BITBUCKET_CLIENT_ID"`
	ClientSecret string   `envconfig:"GATES_BITBUCKET_CLIENT_SECRET"`
	SkipVerity   bool     `envconfig:"GATES_GITLAB_SKIP_VERIFY"`
	Scope        []string `default:"pullrequest:write,account" envconfig:"GATES_BITBUCKET_SCOPE"`
}

// GitLab connection setting
type GitLab struct {
	Server       string   `default:"https://gitlab.com" envconfig:"GATES_GITLAB_SERVER"`
	ClientID     string   `envconfig:"GATES_GITLAB_CLIENT_ID"`
	ClientSecret string   `envconfig:"GATES_GITLAB_CLIENT_SECRET"`
	SkipVerity   bool     `envconfig:"GATES_GITLAB_SKIP_VERIFY"`
	Scope        []string `default:"api,read_user,read_api,read_repository,profile,email" envconfig:"GATES_GITLAB_SCOPE"`
}

// Environ setup configure from environment variables
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
	if c.GitLab.Server != "" {
		providers = append(providers, core.GitLab)
	}
	if c.Bitbucket.Server != "" {
		providers = append(providers, core.Bitbucket)
	}
	return providers
}

// Port opened for the current server
func (server Server) Port() string {
	if server.ServerPort != "" {
		return server.ServerPort
	} else if server.CloudPort != "" {
		return server.CloudPort
	}
	u, err := url.Parse(server.Addr)
	if err != nil {
		return "8080"
	}
	return u.Port()
}

// BaseURL sanitize the Base string to URL format
func (server Server) BaseURL() string {
	base := strings.Trim(server.Base, "/")
	return fmt.Sprintf("/%s", base)
}

// URL of the server
func (server Server) URL() string {
	url := strings.TrimRight(server.Addr, "/")
	url = url + server.BaseURL()
	return strings.TrimRight(url, "/")
}
