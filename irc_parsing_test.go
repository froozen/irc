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
	fmt.Println("Testing:", signal)

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
