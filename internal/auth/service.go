package auth

import (
	pb "github.com/auth_service/api"
	"github.com/auth_service/service/auth"
)

type Layer struct {
	pb.UnimplementedAuthorizationServiceServer
	Auth *auth.Service
}

func NewAuthLayer(auth *auth.Service) *Layer {
	return &Layer{
		Auth: auth,
	}
}
