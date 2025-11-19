module realtime-chat

go 1.22 // Your Go version might be different, that's okay

require (
	github.com/improbable-eng/grpc-web v0.15.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
)

// This replace directive is crucial to fix the dependency conflict
replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240401170217-c3f982113cda

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/rs/cors v1.7.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240325203815-454cdb8f5daa // indirect
	nhooyr.io/websocket v1.8.6 // indirect
)
