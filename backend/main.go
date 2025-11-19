// backend/main.go
package main

import (
    "log"
    "net/http"

    pb "realtime-chat/gen/go" // Change this to your module path

    "github.com/improbable-eng/grpc-web/go/grpcweb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Create a new gRPC server
    grpcServer := grpc.NewServer()
    chatServer := newServer()
    pb.RegisterChatServiceServer(grpcServer, chatServer)

    // Wrap the gRPC server with the gRPC-Web proxy
    wrappedGrpc := grpcweb.WrapServer(grpcServer,
        // Enable CORS
        grpcweb.WithOriginFunc(func(origin string) bool {
            return true
        }),
    )

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if wrappedGrpc.IsAcceptableGrpcCorsRequest(r) || wrappedGrpc.IsGrpcWebRequest(r) {
            wrappedGrpc.ServeHTTP(w, r)
            return
        }
        // Fallback to other handlers
        http.DefaultServeMux.ServeHTTP(w, r)
    })

    log.Println("Starting server on :9090")
    if err := http.ListenAndServe(":9090", handler); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}