package auth

import (
	"github.com/jackc/pgx"
)

type Service struct {
	db *pgx.Conn
}

func NewService(db *pgx.Conn) *Service {
	return &Service{
		db: db,
	}
}
