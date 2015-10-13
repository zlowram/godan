package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Name  string
	Rpcmq Rpcmq
	Monmq Monmq
	Proxy Proxy
}

type Rpcmq struct {
	Host         string
	Port         string
	Queue        string
	Exchange     string
	ExchangeType string
}

type Monmq struct {
	Host     string
	Port     string
	Parallel int
}

type Proxy struct {
	Host     string
	Port     string
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