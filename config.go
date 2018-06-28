package main

import (
	"flag"
	"strings"

	"github.com/pkg/errors"
)

// ConfigType contains configuration information from commandline, etc
type ConfigType struct {
	ListenAddress string
	DialAddresses []string
}

func loadConfig() (ConfigType, error) {
	var config ConfigType
	var dialAddresses string

	flag.StringVar(&config.ListenAddress, "listen", "",
		"address for the server to listen to")
	flag.StringVar(&dialAddresses, "dial", "",
		"colon separated list of addresses for client(s) to dial")

	flag.Parse()

	if config.ListenAddress == "" {
		return ConfigType{}, errors.Errorf("You must specify a listen address with --listen")
	}

	if dialAddresses == "" {
		return ConfigType{}, errors.Errorf("You must specify a colong separated list of dial addresses with --dial")
	}

	config.DialAddresses = strings.Split(dialAddresses, ":")

	return config, nil
}
