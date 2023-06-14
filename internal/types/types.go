package types

import (
	"net"
)

// Connects network agents over TCP.
type Connection net.Conn

// Current context of the application.
type Context struct {
	SubscriptionPort  int
	EventsPort        int
	UsersPool         map[int]Connection
	FollowRegistry    map[int]map[int]bool
	EventQueue        map[int]Event
	EventChannel      chan Event
	DeadLetterChannel chan string
}

// Represents the base format of an event.
type Event struct {
	Sequence       int
	EventType      string
	EmitterUserID  int
	ReceiverUserID int
	Payload        string
}

// Listens to a port and executes Process() after receiving an LR character.
type Server interface {
	Listen(port int)
}

// Processes a server message.
type Handler func(conn Connection, message string)

// Starts processor execution independent of a TCP connection.
type ProcessorStart func(int) int

// Starts processing events as they are added into the event queue.
type Processor interface {
	Start()
}
