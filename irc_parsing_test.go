package irc

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	fmt.Println("--- TestParse")
	if _, err := Parse("JOIN"); err == nil {
		t.Error("Accepts signals with to few arguments.")
	}

	parseTest(t, ":nick!user@host JOIN", "nick", "user", "host", "JOIN")
	parseTest(t, ":nick@host JOIN", "nick", "", "host", "JOIN")
	parseTest(t, ":nick!user JOIN", "nick", "user", "", "JOIN")
	parseTest(t, "JOIN #channel", "", "", "", "JOIN", "#channel")
	parseTest(t, ":node.server.net PING node.server.net", "node.server.net", "",
		"", "PING", "node.server.net")
	parseTest(t, "PRIVMSG #channel :Some message", "", "", "", "PRIVMSG", "#channel", "Some message")
}

func parseTest(t *testing.T, signal, name, user, host, _type string, args ...string) {
	event, err := Parse(signal)
	if err != nil {
		t.Error("Parsing failed:", err)
	}

	if event.Name != name {
		t.Errorf("Parsing name failed: parsed: %s, expected %s", event.Name, name)
	}

	if event.User != user {
		t.Errorf("Parsing user failed: parsed: %s, expected %s", event.User, user)
	}

	if event.Host != host {
		t.Errorf("Parsing host failed: parsed: %s, expected %s", event.Host, host)
	}

	if event.Type != _type {
		t.Errorf("Parsing type failed: parsed: %s, expected %s", event.Type, _type)
	}

	if len(args) != len(event.Arguments) {
		t.Errorf("Parsing arguments failed: parsed: %s, expected %s", event.Arguments, args)
		return
	}

	for i, arg := range args {
		if event.Arguments[i] != arg {
			t.Errorf("Parsing arguments failed: parsed: %s, expected %s", event.Arguments, args)
		}
	}
}

func TestToString(t *testing.T) {
	fmt.Println("--- TestToString")

	join := NewEvent("JOIN")
	join.SetArguments("#channel")
	toStringTest(t, join, "JOIN #channel")

	privmsg := NewEvent("PRIVMSG")
	privmsg.SetArguments("#channel", "Hello there")
	toStringTest(t, privmsg, "PRIVMSG #channel :Hello there")
}

func toStringTest(t *testing.T, event *Event, expected string) {
	if event.String() != expected {
		t.Errorf("Converting failed: converted: %s, expected %s", event, expected)
	}
}

func TestFindOrigin(t *testing.T) {
	fmt.Println("--- TestFindOrigin")

	origin := NewEvent("PRIVMSG")
	origin.SetArguments("#channel", "Hello friend")
	if origin.FindOrigin() != "#channel" {
		t.Error("Incorrect origin: found:", origin.FindOrigin(), "expected: #channel")
	}

	origin.SetArguments("myname", "Hello again")
	origin.Name = "nickname"
	if origin.FindOrigin() != "nickname" {
		t.Error("Incorrect origin: found:", origin.FindOrigin(), "expected: nickname")
	}
}

func ExampleEvent_FindOrigin() {
	origin := NewEvent("PRIVMSG")
	origin.SetArguments("#channel", "Hello friend")
	fmt.Println(origin.FindOrigin())

	origin.SetArguments("myname", "Hello again")
	origin.Name = "nickname"
	fmt.Println(origin.FindOrigin())

	// Output:
	// #channel
	// nickname
}
