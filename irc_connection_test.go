package irc

import (
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	fmt.Println("--- TestConnection")

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

func TestReceive(t *testing.T) {
	fmt.Println("--- TestReceive")

	con, err := Connect("irc.freenode.net", 6667)
	if err != nil {
		t.Fatal(err)
	}
	defer con.Close()

	for i := 0; i < 3; i++ {
		ev, err := con.Receive()
		if err != nil {
			t.Fatal("Error while reading:", err)
		} else if ev.Type != "NOTICE" {
			// The first three events on freenode should be notices
			t.Error(fmt.Sprintln("Unexpected signal:", ev.Type, ev.Arguments))
		}
	}
}

func TestSend(t *testing.T) {
	fmt.Println("--- TestSend")

	con, err := Connect("irc.freenode.net", 6667)
	if err != nil {
		t.Fatal(err)
	}

	nick := NewEvent("NICK")
	nick.SetArguments("testnick")
	con.Send(nick)

	user := NewEvent("USER")
	user.SetArguments("testnick", "localhost", "localhost", "library test")
	con.Send(user)

	for {
		ev, err := con.Receive()
		if err != nil {
			t.Fatal(err)
		}
		// Server recognized the client
		if ev.Type == "433" || ev.Type == "MODE" {
			break
		}
	}
	con.Close()
}
