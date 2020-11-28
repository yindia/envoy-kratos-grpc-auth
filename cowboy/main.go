package main

import (
	"flag"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"

	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/config"

	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/server"
	pb "github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/proto"
)

var configPath = flag.String(
	"config",
	"config.json",
	"config path",
)

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	config.Init(*configPath)

	lis, err := net.Listen("tcp", config.GetGrpcAddress())
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlatformServiceServer(s, server.New())

	// Serve gRPC Server
	log.Info("Serving gRPC on http://", config.GetGrpcAddress())
	go func() {
		log.Fatal(s.Serve(lis))
	}()
	var stopChan chan bool
	<-stopChan
}
