package auth

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	auth_service "github.com/proto_test/api"
)

func (s *Layer) SignUp(ctx context.Context, req *auth_service.SignUpRequest) (response *empty.Empty, err error) {
	err = s.Auth.SignUp(ctx, req)
	if err != nil {
		err = fmt.Errorf("error in SignUp func : %v", err)
		return
	}
	return &empty.Empty{}, nil
}
