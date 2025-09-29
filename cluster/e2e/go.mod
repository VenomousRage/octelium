module github.com/octelium/octelium/cluster/e2e

go 1.24.7

require (
	github.com/octelium/octelium/cluster/common v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/octelium/octelium/apis v0.0.0-00010101000000-000000000000 // indirect
	github.com/octelium/octelium/pkg v0.0.0-00010101000000-000000000000 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.71.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/octelium/octelium/apis => ../../apis

replace github.com/octelium/octelium/pkg => ../../pkg

replace github.com/octelium/octelium/cluster/common => ../common

replace github.com/octelium/octelium/cluster/rscserver => ../rscserver

replace github.com/octelium/octelium/cluster/apiserver => ../apiserver

replace github.com/octelium/octelium/cluster/octovigil => ../octovigil
