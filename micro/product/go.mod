module shopping/product

go 1.13

replace (
	github.com/golang/lint => /Users/admin1/gopath/pkg/mod/golang.org/x/lint@v0.0.0-20200302205851-738671d3881b
	github.com/testcontainers/testcontainer-go => /Users/admin1/gopath/pkg/mod/github.com/testcontainers/testcontainers-go@v0.5.1
)

require (
	github.com/golang/protobuf v1.4.1
	github.com/google/btree v1.0.0 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/micro/go-config v1.1.0
	github.com/micro/go-log v0.1.0
	github.com/micro/go-micro v1.18.0
)
