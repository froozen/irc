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

	con, _ := Connect("irc.freenode.net", 6667)
	for i := 0; i < 3; i++ {
		ev, err := con.Receive()
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

func TestSend(t *testing.T) {
	fmt.Println("--- TestSend")

	con, _ := Connect("irc.freenode.net", 6667)

	nick := &Event{Type: "NICK"}
	nick.Arguments = make([]string, 1)
	nick.Arguments[0] = "testnick"
	con.Send(nick)

	user := &Event{Type: "USER"}
	user.Arguments = make([]string, 4)
	user.Arguments[0] = "testnick"
	user.Arguments[1] = "localhost"
	user.Arguments[2] = "localhost"
	user.Arguments[3] = "library test"
	con.Send(user)

	for {
		ev, _ := con.Receive()
		// Server recognized the client
		if ev.Type == "433" || ev.Type == "MODE" {
			break
		}
	}
	con.Close()
}
