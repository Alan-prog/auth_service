package main

import (
	"context"
	pb "github.com/auth_service/api"
	auth2 "github.com/auth_service/internal/auth"
	"github.com/auth_service/service/auth"
	"github.com/auth_service/tools/db"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

const (
	localServiceAddr    = "127.0.0.1:8090"
	serviceDockerPort   = ":8080"
	serviceLocalPort    = ":8000"
	dockerInsidePortDB  = uint16(5432)
	dockerOutsidePortDB = uint16(6000)

	flagInDocker = "in_docker"
	flagServer   = "in_server"
)

func init() {
	pflag.Bool(flagInDocker, false, "if run inside the docker")
	pflag.Bool(flagServer, false, "if run on server")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("viper.BindPFlags: %s", err)
	}
}

func main() {
	var (
		dbPort      = dockerOutsidePortDB
		proxyPort   = serviceLocalPort
		serviceAddr = localServiceAddr
	)

	ctx := context.Background()

	if viper.GetBool(flagInDocker) {
		dbPort = dockerInsidePortDB
		proxyPort = serviceDockerPort
	}

	dbAdp, err := db.NewDbConnector(ctx, "postgres", "somepass", "localhost", "postgres", dbPort)
	if err != nil {
		log.Fatal("error while connecting to db")
	}
	go runGRPCServer(serviceAddr, dbAdp)
	runHTTPProxy(ctx, serviceAddr, proxyPort)
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

func runHTTPProxy(ctx context.Context, serviceAddr string, proxyAddr string) {
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
