package irc

import (
	"errors"
	"fmt"
	"net"
)

// IRCConnection represents a connection to an IRC server.
type IRCConnection struct {
	connection net.Conn
	address    string
	port       int
}

// Connect connects to an IRC server.
func Connect(address string, port int) (*IRCConnection, error) {
	con := new(IRCConnection)
	con.address = address
	con.port = port

	var err error
	con.connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, errors.New(fmt.Sprintln("While connecting:", err))
	}

	return con, nil
}

// Close closes the IRCConnection.
func (con *IRCConnection) Close() error {
	err := con.connection.Close()
	if err != nil {
		return errors.New(fmt.Sprintln("While closing:", err))
	}
	return nil
}
