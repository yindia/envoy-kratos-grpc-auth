package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"

	auth_v3 "github.com/evalsocket/cloud-factory/services/authz/pkg/auth/v3"
)

func main() {
	port := flag.Int("port", 9001, "gRPC port")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", *port, err)
	}
	if err != nil {
		log.Fatalf("failed to load user data:%v", err)
	}
	gs := grpc.NewServer()

	envoy_service_auth_v3.RegisterAuthorizationServer(gs, auth_v3.New())

	log.Printf("starting gRPC server on: %d\n", *port)

	gs.Serve(lis)
}
