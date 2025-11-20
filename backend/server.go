package main

import (
	"log"
	"sync"
	"time"

	pb "realtime-chat/backend/gen/go"
	// All unused imports have now been removed.
)

// connection represents a single client connection with its stream.
type connection struct {
	stream pb.ChatService_ChatStreamServer
	user   string
	err    chan error
}

// Server holds all the active connections and handles broadcasting.
type Server struct {
	pb.UnimplementedChatServiceServer
	connections []*connection
	mu          sync.Mutex // Mutex to protect access to the connections slice
}

// NewServer creates a new instance of our chat server.
func NewServer() *Server {
	return &Server{}
}

// ChatStream is the main RPC method for the bi-directional stream.
func (s *Server) ChatStream(stream pb.ChatService_ChatStreamServer) error {
	conn := &connection{
		stream: stream,
		err:    make(chan error),
	}

	s.mu.Lock()
	s.connections = append(s.connections, conn)
	s.mu.Unlock()

	log.Println("New user connected.")
	go s.receiveMessages(conn)
	return <-conn.err
}

// receiveMessages runs in a separate goroutine for each client.
func (s *Server) receiveMessages(conn *connection) {
	for {
		msg, err := conn.stream.Recv()
		if err != nil {
			log.Printf("Client %s disconnected: %v", conn.user, err)
			s.removeConnection(conn)
			conn.err <- err // Signal that this connection is done
			return
		}

		if conn.user == "" {
			conn.user = msg.GetMessage().GetUser()
		}

		log.Printf("Received message from %s: %s", conn.user, msg.GetMessage().GetMessage())
		s.broadcast(msg)
	}
}

// broadcast sends a message to all active connections.
func (s *Server) broadcast(msg *pb.StreamMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg.GetMessage().Timestamp = time.Now().Unix()

	for _, conn := range s.connections {
		go func(c *connection) {
			if err := c.stream.Send(msg); err != nil {
				log.Printf("Error sending to %s: %v", c.user, err)
			}
		}(conn)
	}
}

// removeConnection safely removes a connection from the pool and announces the departure.
func (s *Server) removeConnection(conn *connection) {
	s.mu.Lock()

	var activeConnections []*connection
	var departingUser string

	// Create a new slice containing all connections except the one that is leaving.
	for _, c := range s.connections {
		if c != conn {
			activeConnections = append(activeConnections, c)
		} else {
			departingUser = c.user
		}
	}
	s.connections = activeConnections
	
	s.mu.Unlock() // Unlock the mutex before broadcasting

	log.Printf("Connection removed. Active connections: %d", len(s.connections))
	
	// If the user had a name, announce their departure to the remaining users.
	if departingUser != "" {
		leaveMsg := &pb.StreamMessage{
			Event: &pb.StreamMessage_Message{
				Message: &pb.ChatMessage{
					User:    "Server",
					Message: `⬅️ ` + departingUser + ` has left the chat.`,
				},
			},
		}
		// This is the corrected line: using `leaveMsg`
		s.broadcast(leaveMsg)
	}
}