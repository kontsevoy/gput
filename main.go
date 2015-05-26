package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(0)

	// read the command line arguments and the config file
	params, err := ProcessConfig()
	exitIf(err)

	// need to dump config?
	if params.Command == CommandTemplate {
		fmt.Println(ConfigTemplate)
		return
	}

	// authenticate into Rackspace:
	session, err := authenticate(params.ApiKey)
	if err != nil {
		exitIf(fmt.Errorf("%v when trying to authenticate\n", err))
	}
	session.Region = params.Region

	switch params.Command {
	case CommandList:
		if params.Parameter == "" {
			// list containers:
			session.listContainers()
		} else {
			// list objects in a container:
			params.Container = params.Parameter
			session.listObjects(params.Container)
		}
	case CommandPut:
		// upload an object:
		if params.Parameter == "" {
			exitWith("No file specified")
		}
		file, err := os.Open(params.Parameter)
		exitIf(err)

		filename := filepath.Base(params.Parameter)
		if params.SecondParameter != "" {
			filename = params.SecondParameter
		}
		urls := session.upsertObject(file, params.Container, filename, params.TTL)
		printUrls(urls, params.CnameHost)

	case CommandDel:
		if params.Container == "" {
			exitWith("No container specified")
		}
		session.deleteObject(params.Container, params.Parameter)

	default:
		fmt.Println("Unrecoognized command")
	}
}

// Takes a slice of URLs and prints them, replacing their hostname with a given one, if provided
func printUrls(urls []string, hostname string) {
	for _, url := range urls {
		if hostname != "" {
			url = replaceHostnameIn(url, hostname)
		}
		fmt.Println(url)
	}
}
