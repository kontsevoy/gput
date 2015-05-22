package main

import (
	"fmt"
	"strings"
)

func main() {
	// authenticate into Rackspace:
	rax, err := authenticate()
	if err != nil {
		fmt.Printf("%v when trying to authenticate\n", err)
		return
	}

	// Make HTTP put call:
	url := rax.getEntryPoint("DFW", "object-store")
	fmt.Println(url)

	rax.listContainers(url)
	rax.upsertObject(url, strings.NewReader("hello world"), "hello.txt")
}
