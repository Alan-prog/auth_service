package auth

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/proto_test/grpc/api"
)

func (s *Layer) Alive(ctx context.Context, req *empty.Empty) (response *pb.AliveResponse, err error) {
	resp, err := s.Auth.Alive(ctx)
	if err != nil {
		err = fmt.Errorf("error in Alive func : %v", err)
		return
	}
	return &pb.AliveResponse{Message: resp}, nil
}
