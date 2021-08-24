package auth

import (
	pb "github.com/proto_test/grpc/api"
	"github.com/proto_test/service/auth"
)

type Layer struct {
	pb.UnimplementedReverseServer
	Auth *auth.Service
}

func NewAuthLayer(auth *auth.Service) *Layer {
	return &Layer{
		Auth: auth,
	}
}
