package auth

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/auth_service/api"
)

func (s *Layer) ApprovePhoneNumber(ctx context.Context, req *pb.ApproveRequest) (response *empty.Empty, err error) {
	err = s.Auth.ApprovePhoneNumber(ctx)
	if err != nil {
		err = fmt.Errorf("error in ApprovePhoneNumber func : %v", err)
		return
	}
	return &empty.Empty{}, nil
}
