package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server interface {
	HandleConnection(conn net.Conn)
	Setup(port string)
}

// MessageHandler is a function that processes a message and returns a response string or an error.
type MessageHandler func(msg string) (string, error)

type TCPServer struct {
	Handler MessageHandler
}

func (s TCPServer) Setup(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

// HandleConnection handles the lifecycle of a single network connection.
func (s TCPServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %v", err)
			return
		}
		msg := strings.TrimSpace(message)

		response, err := s.Handler(msg)
		if err != nil {
			fmt.Fprintf(conn, "Error process your command %s\n", err.Error())
			continue
		}
		conn.Write([]byte(response + "\n"))
	}
}