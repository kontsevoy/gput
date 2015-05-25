package main

/*
This file contains code to deal with parameters/configuration, including:
	- command line arguments
	- ini file
*/

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

type Params struct {
	ConfigPath      string // path to a config file
	ApiKey          string
	ApiUser         string
	Container       string
	ObjectName      string
	Command         string
	Region          string
	Parameter       string
	SecondParameter string
	TTL             int // time to live (for file uploads). 0 means "forever"
}

const (
	DefaultConfigFile = "~/.gput.ini"
	DefaultCommand    = "put"
)

const ConfigTemplate = `; Save this into ~/.gput.ini
[Auth]
key=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
username=xxxxxxxxxxxxx

; Section to configure default cloud files options
[Cloud Files]
container=public-container
region=DFW
`

// list of possible commands:
const (
	CommandPut      = "put"    // uploads object into a container
	CommandList     = "list"   // lists containers or/and objects
	CommandDel      = "delete" // deletes object/container
	CommandTemplate = "gen"    // generates a dumps a sample config file
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
	if params.Region == "" {
		params.Region = iniConf.Get("CloudFiles", "region")
	}
	params.Region = strings.ToUpper(params.Region)

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
	if p.Region == "" {
		return errors.New("Region (like 'DFW', 'ORD', etc) is not specified")
	}
	if p.Command == "" {
		p.Command = DefaultCommand
	}

	return nil
}

func parseCommandLine() (params Params) {
	flag.StringVar(&params.ConfigPath, "config", "", "path to a config file")
	flag.StringVar(&params.ApiKey, "key", "", "Rackspace API key")
	flag.StringVar(&params.ApiUser, "user", "", "Rackspace API username")
	flag.StringVar(&params.Container, "container", "", "Cloud Files container name")
	flag.StringVar(&params.Region, "region", "", "Cloud Files region to use, like DFW, ORD, etc")
	flag.IntVar(&params.TTL, "ttl", 0, "Time to live in seconds. 0 (default) means forever")

	flag.Usage = func() {
		fmt.Printf("gput is a client for Rackspace Cloud files\n\n")
		fmt.Printf("Usage:\n\tgput <options> <command> <file>\n\n")
		fmt.Printf("Commands:\n")
		fmt.Printf("\tput   :\tUpload file into a container. This command executes by default.\n")
		fmt.Printf("\tlist  :\tLlist containers or files within a container\n")
		fmt.Printf("\tdelete:\tDelete file in a container\n")
		fmt.Printf("\tgen   :\tGenerate a template config file\n\n")
		fmt.Printf("Options:\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("\t-%v : %v\n", f.Name, f.Usage)
		})
		fmt.Printf("\nExamples:\n")
		fmt.Printf("\tgput -region DFW list\n")
		fmt.Printf("\tgput -region DFW -container public-container list\n")
		fmt.Printf("\tgput -region DFW list public-container\n")
		fmt.Printf("\tgput -region DFW -container public-container delete example.txt\n")
		fmt.Printf("\tgput -region DFW put example.txt\n")
		fmt.Printf("\tgput -ttl 600 put example.txt\n")
	}

	flag.Parse()

	params.Command = flag.Arg(0)
	params.Parameter = flag.Arg(1)
	params.SecondParameter = flag.Arg(2)

	// first command is a file name? assume it's a parameter
	if fileExists(params.Command) {
		params.Command = DefaultCommand
		params.Parameter = flag.Arg(0)
		params.SecondParameter = flag.Arg(1)
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
