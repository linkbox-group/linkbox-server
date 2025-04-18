package model

import "github.com/golang-jwt/jwt/v4"

type TokenType int

const (
	TokenTypeAccess TokenType = iota + 1
	TokenTypeRefresh
)

type TokenClaims struct {
	UserId string    `json:"userId"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}
