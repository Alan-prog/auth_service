package auth

import "context"

func (s *Service) Alive(ctx context.Context) (response string, err error) {
	return "It is okay", nil
}
