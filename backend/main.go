package main

import (
	"log"
	"net/http"
	"net" // You need this import for net.Listen

	pb "realtime-chat/gen/go"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	// The "google.golang.org/grpc/credentials/insecure" line has been removed
)

func main() {
	// We will now listen on two different ports.
	// :9090 for the pure gRPC server
	// :8080 for the gRPC-Web proxy that the browser will connect to
	
	// Create a TCP listener for the main gRPC server
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen on port 9090: %v", err)
	}

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()
	
	// Create an instance of our chat server implementation
	chatServer := NewServer()

	// Register the chat server with the gRPC server
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	log.Println("gRPC server listening on :9090")

	// Start the pure gRPC server in a separate goroutine
	// This is for potential non-web clients in the future
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	
	// Wrap the gRPC server with the gRPC-Web proxy
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		// Enable CORS to allow requests from our frontend's origin
		grpcweb.WithOriginFunc(func(origin string) bool {
			// For development, allowing all origins is fine.
			return true
		}),
	)
	
	// Create a handler for the gRPC-Web proxy
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedGrpc.ServeHTTP(w, r)
	})

	log.Println("gRPC-Web proxy listening on :8080")

	// Start the HTTP server for the gRPC-Web proxy on port 8080
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to serve gRPC-Web: %v", err)
	}
}