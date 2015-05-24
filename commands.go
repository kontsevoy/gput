package main

import (
	"errors"
	"flag"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

type Params struct {
	ConfigPath string
	ApiKey     string
	ApiUser    string
	Container  string
	ObjectName string
	Command    string
	Parameter  string
}

const (
	DefaultConfigFile = "~/.gput.ini"
	DefaultCommand    = "put"
)

// list of possible commands:
const (
	CommandPut  = "put"    // uploads object into a container
	CommandList = "list"   // lists containers or/and objects
	CommandDel  = "delete" // deletes object/container

)

var (
	PossibleCommands = []string{CommandPut, CommandList, CommandDel}
)

// Parses the command line arguments and reads the config file if specified
func ProcessConfig() (params Params, err error) {
	var iniConf IniConfig
	// get configs from both the command line and the config file:
	params = parseCommandLine()
	iniConf, err = parseConfigFile(&params)
	if err != nil {
		return
	}

	// merge ini and command line values (command line overrides):
	if params.ApiKey == "" {
		params.ApiKey = iniConf.Get("Auth", "key")
	}
	if params.ApiUser == "" {
		params.ApiUser = iniConf.Get("Auth", "username")
	}
	if params.Container == "" {
		params.Container = iniConf.Get("CloudFiles", "container")
	}

	// check correctness:
	err = checkConfig(&params)
	return
}

// Checks all configuration settings for sanity
func checkConfig(p *Params) error {
	if p.ApiKey == "" {
		return errors.New("API key is required")
	}
	if p.ApiUser == "" {
		return errors.New("API username is required")
	}
	if p.Container == "" {
		return errors.New("Container is not specified")
	}
	if p.Command == "" {
		p.Command = DefaultCommand
	}

	return nil
}

func parseCommandLine() (params Params) {
	flag.StringVar(&params.ConfigPath, "config", "", "Path to the config file")
	flag.StringVar(&params.ApiKey, "key", "", "Rackspace API key")
	flag.StringVar(&params.ApiUser, "user", "", "Rackspace API username")
	flag.StringVar(&params.Container, "container", "", "Default Cloud Files container name")
	flag.Parse()

	params.Command = flag.Arg(0)
	params.Parameter = flag.Arg(1)

	// first command is a file name? assume it's a parameter
	if fileExists(params.Command) {
		params.Command = DefaultCommand
		params.Parameter = flag.Arg(0)
	}

	return
}

// parseConfigFile() tries to open the command-line specified config file, and
// if it does not exist, goes for the default one.
func parseConfigFile(p *Params) (iniConf IniConfig, err error) {
	// check if the command line specified config file exists
	if p.ConfigPath != "" {
		p.ConfigPath, _ = filepath.Abs(p.ConfigPath)
		if !fileExists(p.ConfigPath) {
			err = errors.New("Configuration file does not exist: " + p.ConfigPath)
			return
		}
	} else {
		p.ConfigPath = DefaultConfigFile

		// expand ~/ into full path:
		if me, e := user.Current(); e == nil {
			p.ConfigPath = strings.Replace(p.ConfigPath, "~/", me.HomeDir+"/", 1)
		}
		if !fileExists(p.ConfigPath) {
			log.Printf("No such file: %v Using only command line args", p.ConfigPath)
			return
		}
	}
	// parse config files:
	return ParseIniFile(p.ConfigPath)
}
