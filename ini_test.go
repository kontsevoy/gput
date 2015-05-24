package main

import (
	"fmt"
	"testing"
)

func TestIni(t *testing.T) {
	conf, err := ParseIniFile("fixtures/gput.conf")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Key: %v\n", conf.Get("Auth", "key"))
	fmt.Println(conf)
}
