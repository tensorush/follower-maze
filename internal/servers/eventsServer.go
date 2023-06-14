package servers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"follower-maze/internal/types"
)

type eventsServer struct {
	handler types.Handler
}

// Returns a server instance that closes the connection whenever it stops receiving messages.
func NewEventsServer(handler types.Handler) types.Server {
	return eventsServer{
		handler: handler,
	}
}

// Starts the server on the given port.
func (s eventsServer) Listen(port int) {
	eventListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer eventListener.Close()

	fmt.Printf("Listening for operators on %d\n", port)

outer:
	for {
		conn, err := eventListener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(conn)

		for {
			payloadRaw, err := reader.ReadString('\n')

			if err == io.EOF {
				conn.Close()
				continue outer
			} else if err != nil {
				log.Fatal(err)
			}

			payload := strings.TrimSpace(payloadRaw)

			fmt.Printf("Message received: %s\n", payload)

			s.handler(conn, payload)
		}
	}
}
