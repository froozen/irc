// irc provides an implementation of the Internet Relay Chat protocol.
// On top of that, it offers some shortcuts for easier and more comfortable
// usage.
//
// The implementation tries to comply to this specification:
// http://tools.ietf.org/html/rfc1459
package irc

import (
	"errors"
	"fmt"
	"strings"
)

// IRCEvent represents a parsed IRC signal
type IRCEvent struct {
	// Data representing the signal origin
	Name string
	User string
	Host string

	// Data representing the contents of the signal
	Type      string
	Arguments []string
}

// Parse parses an IRC signal.
func Parse(signal string) (*IRCEvent, error) {
	var args []string
	signalArgs := strings.Split(signal, " ")

	// The next parsing steps definetly won't fail if this passes
	if len(signalArgs) < 2 {
		return nil, errors.New(fmt.Sprintf("\"%s\": Not enough arguments."))
	}

	ircEvent := new(IRCEvent)

	// Signal has source
	if strings.HasPrefix(signalArgs[0], ":") {
		// Structure :nick!user@host (only nick is guaranteed, though)
		source := strings.TrimPrefix(signalArgs[0], ":")

		// Source contains user
		if strings.Contains(source, "!") {
			split := strings.Split(source, "!")
			// Source contains host as well
			if strings.Contains(split[1], "@") {
				spilt2 := strings.Split(split[1], "@")
				ircEvent.User = spilt2[0]
				ircEvent.Host = spilt2[1]
				source = split[0]
			} else {
				ircEvent.User = split[1]
				source = split[0]
			}
		}

		// Source contains host
		if strings.Contains(source, "@") {
			split := strings.Split(source, "@")
			ircEvent.Host = split[1]
			source = split[0]
		}

		// At this point, only the name will remain in source
		ircEvent.Name = source

		// This means the second argument is the type
		ircEvent.Type = signalArgs[1]

		args = signalArgs[2:]
	} else {
		ircEvent.Type = signalArgs[0]

		args = signalArgs[1:]
	}

	// Parse args
	ircEvent.Arguments = args
	for i, arg := range args {
		// Detect the trailing argument and join it
		if strings.HasPrefix(arg, ":") {
			ircEvent.Arguments = args[:i+1]
			ircEvent.Arguments[i] = strings.Join(args[i:], " ")
			ircEvent.Arguments[i] = strings.TrimPrefix(ircEvent.Arguments[i], ":")
			break
		}
	}

	return ircEvent, nil
}
