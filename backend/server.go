// backend/server.go
package main

import (
    "log"
    "net"
    "sync"
    "time"

    pb "realtime-chat/gen/go" // Change this to your module path

    "google.golang.org/grpc"
)

type connection struct {
    stream pb.ChatService_ChatStreamServer
    user   string
    err    chan error
}

type Server struct {
    pb.UnimplementedChatServiceServer
    connections []*connection
    mu          sync.Mutex
}

func (s *Server) ChatStream(stream pb.ChatService_ChatStreamServer) error {
    conn := &connection{
        stream: stream,
        err:    make(chan error),
    }

    s.mu.Lock()
    s.connections = append(s.connections, conn)
    s.mu.Unlock()

    log.Println("New user connected")

    // Goroutine to receive messages
    go func() {
        for {
            msg, err := stream.Recv()
            if err != nil {
                log.Printf("Error receiving message from client: %v", err)
                conn.err <- err
                return
            }

            // First message is the username
            if conn.user == "" {
                conn.user = msg.GetMessage().GetUser()
            }

            log.Printf("Received message: %s from %s", msg.GetMessage().GetMessage(), conn.user)
            s.broadcast(msg)
        }
    }()

    return <-conn.err
}

func (s *Server) broadcast(msg *pb.StreamMessage) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Add server timestamp
    msg.GetMessage().Timestamp = time.Now().Unix()

    for _, conn := range s.connections {
        if err := conn.stream.Send(msg); err != nil {
            log.Printf("Error sending message to user %s: %v", conn.user, err)
        }
    }
}

func newServer() *Server {
    return &Server{}
}