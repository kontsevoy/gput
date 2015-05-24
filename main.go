package main

import (
	"fmt"
	"strings"
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

	session.listContainers()
	session.listObjects(params.Container)
	session.upsertObject(strings.NewReader("hello world"), params.Container, "hello.txt")
}
