package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "realtime-chat/backend/gen/go"

	"google.golang.org/protobuf/types/known/emptypb"
)

// connection represents a single client connection with its stream.
type connection struct {
	stream pb.ChatService_SubscribeMessagesServer
	user   string
	active bool
	error  chan error
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

// SendMessage is a unary RPC to send a message to the server.
func (s *Server) SendMessage(ctx context.Context, msg *pb.ChatMessage) (*emptypb.Empty, error) {
	log.Printf("Received message from %s: %s", msg.User, msg.Message)
	
	// Broadcast the message to all subscribers
	streamMsg := &pb.StreamMessage{
		Event: &pb.StreamMessage_Message{
			Message: msg,
		},
	}
	s.broadcast(streamMsg)
	
	return &emptypb.Empty{}, nil
}

// SubscribeMessages is a server-streaming RPC to receive messages.
func (s *Server) SubscribeMessages(req *pb.SubscribeRequest, stream pb.ChatService_SubscribeMessagesServer) error {
	conn := &connection{
		stream: stream,
		user:   req.User,
		active: true,
		error:  make(chan error),
	}

	s.mu.Lock()
	s.connections = append(s.connections, conn)
	s.mu.Unlock()

	log.Printf("User connected: %s", conn.user)
	
	// Announce user joined
	joinMsg := &pb.StreamMessage{
		Event: &pb.StreamMessage_Message{
			Message: &pb.ChatMessage{
				User:    "Server",
				Message: "➡️ " + conn.user + " has joined the chat.",
				Timestamp: time.Now().Unix(),
			},
		},
	}
	s.broadcast(joinMsg)

	// Keep the stream alive until context is cancelled or error occurs
	<-stream.Context().Done()
	
	log.Printf("User disconnected: %s", conn.user)
	s.removeConnection(conn)
	return nil
}

// broadcast sends a message to all active connections.
func (s *Server) broadcast(msg *pb.StreamMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure timestamp is set if it's a chat message
	if m, ok := msg.Event.(*pb.StreamMessage_Message); ok {
		if m.Message.Timestamp == 0 {
			m.Message.Timestamp = time.Now().Unix()
		}
	}

	for _, conn := range s.connections {
		if !conn.active {
			continue
		}
		go func(c *connection) {
			if err := c.stream.Send(msg); err != nil {
				log.Printf("Error sending to %s: %v", c.user, err)
				// We don't remove here immediately to avoid race conditions on the slice,
				// but the context cancellation in SubscribeMessages will handle cleanup.
			}
		}(conn)
	}
}

// removeConnection safely removes a connection from the pool and announces the departure.
func (s *Server) removeConnection(conn *connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn.active = false
	
	var activeConnections []*connection
	for _, c := range s.connections {
		if c != conn {
			activeConnections = append(activeConnections, c)
		}
	}
	s.connections = activeConnections

	log.Printf("Connection removed. Active connections: %d", len(s.connections))
	
	// Announce departure
	leaveMsg := &pb.StreamMessage{
		Event: &pb.StreamMessage_Message{
			Message: &pb.ChatMessage{
				User:    "Server",
				Message: "⬅️ " + conn.user + " has left the chat.",
				Timestamp: time.Now().Unix(),
			},
		},
	}
	
	// Broadcast to remaining users
	// We need to release the lock to broadcast to avoid deadlock if broadcast tries to lock
	// But broadcast locks too. 
	// To avoid deadlock, we can spawn a goroutine or handle it carefully.
	// Since we are holding the lock, we shouldn't call broadcast directly if it locks.
	// Let's copy the connections and send outside the lock or modify broadcast to not lock?
	// Or just iterate here.
	
	for _, c := range s.connections {
		go func(dest *connection) {
			dest.stream.Send(leaveMsg)
		}(c)
	}
}
