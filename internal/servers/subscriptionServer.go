package servers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"follower-maze/internal/types"
)

type subscriptionServer struct {
	handler types.Handler
}

// Returns an instance of a server.
func NewSubscriptionServer(handler types.Handler) types.Server {
	return subscriptionServer{
		handler: handler,
	}
}

// Starts the server on the given port.
func (s subscriptionServer) Listen(port int) {
	eventListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer eventListener.Close()

	fmt.Printf("Listening for client requests on %d\n", port)

	for {
		conn, err := eventListener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(conn)

		payloadRaw, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		payload := strings.TrimSpace(payloadRaw)

		s.handler(conn, payload)
	}
}
