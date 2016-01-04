package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Local Local
	Rpcmq Rpcmq
	Monmq Monmq
	DB    Database
}

type Local struct {
	Host string
	Port string
}

type Rpcmq struct {
	Host         string
	Port         string
	MsgQueue     string
	ReplyQueue   string
	Exchange     string
	ExchangeType string
}

type Monmq struct {
	Host       string
	Port       string
	ReplyQueue string
	Exchange   string
}

type Database struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
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
