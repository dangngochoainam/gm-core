module github.com/dangngochoainam/gm-core

go 1.24.0

replace github.com/dangngochoainam/gopkg => ../gopkg

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1
	github.com/dangngochoainam/gopkg v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	google.golang.org/genproto/googleapis/api v0.0.0-20260120221211-b8f7ae30c516
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	buf.build/go/protovalidate v1.2.0 // indirect
	cel.dev/expr v0.25.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/google/cel-go v0.28.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3 // indirect
	golang.org/x/exp v0.0.0-20250813145105-42675adae3e6 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)
