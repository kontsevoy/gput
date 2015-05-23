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
	rax, err := authenticate(params.ApiKey)
	if err != nil {
		exitIf(fmt.Errorf("%v when trying to authenticate\n", err))
	}

	// Make HTTP put call:
	url := rax.getEntryPoint("DFW", "object-store")
	fmt.Println(url)

	rax.listContainers(url)
	rax.upsertObject(url, strings.NewReader("hello world"), "hello.txt")
}
