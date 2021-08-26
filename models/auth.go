package models

import (
	"github.com/dgrijalva/jwt-go"
)

var JwtSigningKey = []byte("secret")

type ClaimWithID struct {
	ID string `json:"custom_id"`
	jwt.StandardClaims
}
