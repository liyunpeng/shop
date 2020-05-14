module shopping/notification

go 1.13

replace (
	github.com/golang/lint => /Users/admin1/gopath/pkg/mod/golang.org/x/lint@v0.0.0-20200302205851-738671d3881b
	github.com/testcontainers/testcontainer-go => /Users/admin1/gopath/pkg/mod/github.com/testcontainers/testcontainers-go@v0.5.1
	product => ../product
	user => ../user
)

require (
	github.com/golang/protobuf v1.4.1
	github.com/micro/go-config v1.1.0
	github.com/micro/go-log v0.1.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/broker/grpc v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/broker/rabbitmq v0.0.0-20200119172437-4fe21aa238fd
	product v0.0.0-00010101000000-000000000000
	user v0.0.0-00010101000000-000000000000
)
