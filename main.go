package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// read the command line arguments and the config file
	params, err := ProcessConfig()
	exitIf(err)

	// authenticate into Rackspace:
	session, err := authenticate(params.ApiKey)
	if err != nil {
		exitIf(fmt.Errorf("%v when trying to authenticate\n", err))
	}

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
			log.Fatalf("No file specified")
		}
		file, err := os.Open(params.Parameter)
		if err != nil {
			log.Fatal(err)
		}
		session.upsertObject(file, params.Container, filepath.Base(params.Parameter))

	case CommandDel:
		if params.Container == "" {
			log.Fatal("No container specified")
		}
		session.deleteObject(params.Container, params.Parameter)

	default:
		fmt.Println("No command specified")
	}
}
