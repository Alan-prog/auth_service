package main

import (
	"context"
	"github.com/auth_service/tools/db"
	"github.com/jackc/pgx"
	"log"
	"net"
	"net/http"

	pb "github.com/auth_service/api"
	auth2 "github.com/auth_service/internal/auth"
	"github.com/auth_service/service/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const (
	serviceAddr = "127.0.0.1:8081"
	proxyAddr   = ":8080"
)

func main() {
	ctx := context.Background()

	dbAdp, err := db.NewDbConnector(ctx, "postgres", "somepass", "localhost", "postgres", 5432)
	if err != nil {
		log.Fatal("error while connecting to db")
	}

	go runGRPCServer(serviceAddr, dbAdp)
	runHTTPProxy(ctx)
}

func runGRPCServer(address string, db *pgx.Conn) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gRPCServer := grpc.NewServer()

	// can add here services as much you need
	authSrv := auth.NewService(db)
	authLay := auth2.NewAuthLayer(authSrv)
	pb.RegisterAuthorizationServiceServer(gRPCServer, authLay)

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

	err = pb.RegisterAuthorizationServiceHandler(ctx, grpcGWMUX, grpcConn)
	if err != nil {
		log.Fatalln("error while trying to register server")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcGWMUX)

	log.Printf("starting http server at: %s", proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, mux))
}
