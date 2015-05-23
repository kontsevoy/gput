package main

import "testing"

func TestIni(t *testing.T) {
	err := ParseIniFile("fixtures/gput.conf")
	if err != nil {
		t.Error(err)
	}
}
