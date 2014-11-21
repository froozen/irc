package irc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Connection represents a connection to an IRC server.
type Connection struct {
	connection net.Conn
	address    string
	port       int
	// Scanner for reading
	scanner *bufio.Scanner
	// Communation channels for send routine
	send chan<- *Event
	quit chan<- bool
	// There can be only one receiver at a time
	receiveLock sync.Mutex
}

// Connect connects to an IRC server.
func Connect(address string, port int) (*Connection, error) {
	con := new(Connection)
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

// Close closes the Connection.
func (con *Connection) Close() error {
	con.quit <- true
	err := con.connection.Close()
	if err != nil {
		return errors.New(fmt.Sprintln("While closing:", err))
	}
	return nil
}

// Receive reads and parses the next line received by the Connection
// This call is locked by a mutex, so only one goroutine can receive at a time.
func (con *Connection) Receive() (*Event, error) {
	// Make sure only one routine is receiving at a time
	con.receiveLock.Lock()
	defer con.receiveLock.Unlock()

	if con.scanner.Scan() {
		event, err := Parse(con.scanner.Text())
		if err != nil {
			return nil, errors.New(fmt.Sprintln("While parsing:", err))
		}

		return event, nil

	}

	if con.scanner.Err() != nil {
		return nil, errors.New(fmt.Sprintln("While reading:", con.scanner.Err()))
	}
	return nil, errors.New("Connection closed.")
}

// startSendRoutine starts a routine that schedules the sending of signals to
// prevent flooding
func (con *Connection) startSendRoutine() (chan *Event, chan bool) {
	send := make(chan *Event)
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

// Send sends a signal corresponding to an Event to the server
func (con *Connection) Send(event *Event) {
	con.send <- event
}
