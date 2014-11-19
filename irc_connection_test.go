package irc

import (
	"fmt"
	"testing"
)

func TestIRCConnection(t *testing.T) {
	fmt.Println("--- TestIRCConnection")

	// This is assumed to be invalid
	_, err := Connect("asdfasdfaasdf", 1)
	if err == nil {
		t.Error("Connects to invalid host")
	}

	// This is assumed to be valid
	_, err = Connect("irc.freenode.net", 6667)
	if err != nil {
		t.Error("Doesn't connect to valide host.")
	}
}
