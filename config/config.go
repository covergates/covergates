package config

import "log"

var config *Config

func GetConfig() *Config {
	if config == nil {
		log.Fatal("Config not initialize")
	}
	return config
}

type Config struct {
	Gitea  Gitea
	Github Github
}

type Gitea struct {
	Server       string
	ClientID     string
	ClientSecret string
	SkipVerity   bool
	Scope        []string `default:"repo,repo:status,user:email,read:org"`
}

type Github struct {
	Server       string `default:"https://github.com"`
	APIServer    string
	ClientID     string
	ClientSecret string
	SkipVerity   bool
	Scope        []string `default:"repo,repo:status,user:email,read:org"`
}
