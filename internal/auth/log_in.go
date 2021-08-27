package auth

import (
	"context"
	"fmt"
	auth_service "github.com/auth_service/api"
)

func (s *Layer) LogIn(ctx context.Context, req *auth_service.LogInRequest) (response *auth_service.LogInResponse, err error) {
	resp, err := s.Auth.LogIn(ctx, req)
	if err != nil {
		err = fmt.Errorf("error in LogIn func : %v", err)
		return
	}
	return &resp, nil
}
