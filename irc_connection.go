package irc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

// IRCConnection represents a connection to an IRC server.
type IRCConnection struct {
	connection net.Conn
	address    string
	port       int
	// Scanner for reading
	scanner *bufio.Scanner
	// Communation channels for send routine
	send chan<- *IRCEvent
	quit chan<- bool
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
	con.scanner = bufio.NewScanner(con.connection)
	con.send, con.quit = con.startSendRoutine()

	return con, nil
}

// Close closes the IRCConnection.
func (con *IRCConnection) Close() error {
	con.quit <- true
	err := con.connection.Close()
	if err != nil {
		return errors.New(fmt.Sprintln("While closing:", err))
	}
	return nil
}

// ReadIRCEvent reads and parses the next line received by the IRCConnection
func (con *IRCConnection) ReadIRCEvent() (*IRCEvent, error) {
	if con.scanner.Scan() {
		event, err := Parse(con.scanner.Text())
		if err != nil {
			return nil, errors.New(fmt.Sprintln("While parsing:", err))
		}

		return event, nil

	} else {
		if con.scanner.Err() != nil {
			return nil, errors.New(fmt.Sprintln("While reading:", con.scanner.Err()))
		} else {
			return nil, errors.New("Connection closed.")
		}
	}
}

// startSendRoutine starts a routine that schedules the sending of signals to
// prevent flooding
func (con *IRCConnection) startSendRoutine() (chan *IRCEvent, chan bool) {
	send := make(chan *IRCEvent)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case event := <-send:
				fmt.Fprintf(con.connection, "%s\n\r", event)
				time.Sleep(10 * time.Millisecond)
			case <-quit:
				break
			}
		}
	}()
	return send, quit
}

// SendIRCEvent sends a signal corresponding to an IRCEvent to the server
func (con *IRCConnection) SendIRCEvent(event *IRCEvent) {
	con.send <- event
}
