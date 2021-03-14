package main

import (
	_ "embed"

	"github.com/BurntSushi/toml"
)

type config struct {
	PasswordHash string

	SecurityGroupID string
	Region          string

	Title       string
	Description string
}

//go:embed config.toml
var embeddedConfig string

var currentConfig config

func loadConfig() {
	_, err := toml.Decode(embeddedConfig, &currentConfig)
	if err != nil {
		panic(err)
	}
}
