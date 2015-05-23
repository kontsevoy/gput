package main

import (
	"errors"
	"flag"
	_ "fmt"
)

type Params struct {
	ConfigPath string
	ApiKey     string
	Container  string
}

// Parses the command line arguments and reads the config file if specified
func ProcessConfig() (params Params, err error) {
	params = parseCommandLine()
	err = parseConfigFile(&params)
	if err != nil {
		return
	}
	err = checkConfig(&params)
	return
}

// Checks all configuration settings for sanity
func checkConfig(p *Params) error {
	if p.ApiKey == "" {
		return errors.New("API key is required")
	}
	return nil
}

func parseCommandLine() (params Params) {
	flag.StringVar(&params.ConfigPath, "config", "", "Path to the config file")
	flag.StringVar(&params.ApiKey, "key", "", "Rackspace API key")
	flag.StringVar(&params.Container, "container", "", "Default Cloud Files container name")
	flag.Parse()
	return
}

func parseConfigFile(p *Params) (err error) {
	// no config file? return!
	if p.ConfigPath != "" {
		return nil
	}

	return
}
