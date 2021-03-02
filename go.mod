module github.com/SkyAPM/go2sky

go 1.12

require (
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20200222125558-5a598a2470a0 // indirect
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894 // indirect
	golang.org/x/text v0.3.1-0.20181010134911-4d1c5fb19474 // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	skywalking/network v1.0.0
)

replace skywalking/network v1.0.0 => ./protocol/gen-codes/skywalking/network
