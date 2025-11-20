package main

import (
	"log"
	"net"
	"net/http"

	pb "realtime-chat/backend/gen/go"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // 1. IMPORT THE REFLECTION PACKAGE
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen on port 9090: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatServer := NewServer()
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	// 2. REGISTER THE REFLECTION SERVICE
	// This line allows tools like grpcurl to query for available RPCs.
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :9090")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedGrpc.ServeHTTP(w, r)
	})

	log.Println("gRPC-Web proxy listening on :8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to serve gRPC-Web: %v", err)
	}
}