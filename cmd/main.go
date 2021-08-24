package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/proto_test/grpc/api"
	auth2 "github.com/proto_test/internal/auth"
	"github.com/proto_test/service/auth"
	"google.golang.org/grpc"
)

const (
	serviceAddr = "127.0.0.1:8081"
	proxyAddr   = ":8080"
)

func main() {
	ctx := context.Background()

	go runGRPCServer(serviceAddr)
	runHTTPProxy(ctx)
}

func runGRPCServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gRPCServer := grpc.NewServer()

	// can add here services as much you need
	authSrv := auth.NewService()
	authLay := auth2.NewAuthLayer(authSrv)
	pb.RegisterReverseServer(gRPCServer, authLay)

	log.Printf("starting gRPC service at :%s", address)
	log.Fatal(gRPCServer.Serve(lis))
}

func runHTTPProxy(ctx context.Context) {
	grpcGWMUX := runtime.NewServeMux()
	grpcConn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("error while creating gRPC conn")
	}
	defer grpcConn.Close()

	err = pb.RegisterReverseHandler(ctx, grpcGWMUX, grpcConn)
	if err != nil {
		log.Fatalln("error while trying to register server")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcGWMUX)

	log.Printf("starting http server at: %s", proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, mux))
}
