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

func TestIRCRead(t *testing.T) {
	fmt.Println("--- TestIRCRead")

	con, _ := Connect("irc.freenode.net", 6667)
	for i := 0; i < 3; i++ {
		ev, err := con.ReadIRCEvent()
		if err != nil {
			t.Error("Error while reading:", err)

			// The first three events on freenode should be notices
		} else if ev.Type != "NOTICE" {
			t.Error(fmt.Sprintln("Unexpected signal:", ev.Type, ev.Arguments))
		} else {
			fmt.Println(ev.Type, ev.Arguments)
		}
	}
	con.Close()
}
