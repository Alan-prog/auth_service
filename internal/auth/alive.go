package auth

import (
	"context"
	"fmt"

	pb "github.com/auth_service/api"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Layer) Alive(ctx context.Context, req *empty.Empty) (response *pb.AliveResponse, err error) {
	resp, err := s.Auth.Alive(ctx)
	if err != nil {
		err = fmt.Errorf("error in Alive func : %v", err)
		return
	}
	return &pb.AliveResponse{Message: resp}, nil
}
