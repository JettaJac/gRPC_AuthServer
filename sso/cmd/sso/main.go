package main

// protoc -I proto --go_out=plugins=grpc:cmd proto/sso/sso.proto
// protoc -I proto proto/sso/sso.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
// go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

import (
	"app/internal/config"
	"fmt"
)

func main() {
	fmt.Println("Start app")
	cfg := config.MustLoad()
	fmt.Println(cfg)
}

func setupLogger(env string)
