package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Local       Local
	DB          Database
	Godan       Godan
	DefaultUser DefaultUser
}

type Local struct {
	Host       string
	Port       string
	PrivateKey string
	PublicKey  string
}

type Database struct {
	Host string
	Port string
}

type Godan struct {
	Host string
	Port string
}

type DefaultUser struct {
	Username string
	Password string
}

func loadConfig(configFile string) (Config, error) {
	var config Config

	if _, err := os.Stat(configFile); err != nil {
		return config, err
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return config, err
	}
	return config, nil
}
