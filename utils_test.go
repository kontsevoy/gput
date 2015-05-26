package main

import (
	"testing"
)

func TestUtils(t *testing.T) {
	if "http://newhost.com/v2/CALL?query=value" !=
		replaceHostnameIn("http://example.com/v2/CALL?query=value", "newhost.com") {
		t.Error("replaceHostnameIn() failed")
	}
	if "http://example.com/v2" != replaceHostnameIn("http://example.com/v2", "") {
		t.Error("replaceHostnameIn() should ignore empty string")

	}
}
