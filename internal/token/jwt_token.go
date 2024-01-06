package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct {
	secret   string
	duration time.Duration
}

func NewJWTMaker(secret string, duration time.Duration) *JWTMaker {
	return &JWTMaker{secret, duration}
}

func (maker *JWTMaker) GenerateToken(username string, role string) (string, *UserPayload, error) {
	claim, err := NewPayload(username, role, maker.duration)
	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := jwtToken.SignedString([]byte(maker.secret))
	return token, claim, err
}

func (maker *JWTMaker) VerifyToken(token string) (*UserPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidToken
		}
		return []byte(maker.secret), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &UserPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ExpiredToken) {
			return nil, ExpiredToken
		}
	}

	payload, ok := jwtToken.Claims.(*UserPayload)
	if !ok {
		return nil, InvalidToken
	}

	return payload, nil
}

func (maker *JWTMaker) GetUsernameFromToken(token string) (string, error) {
	if claims, ok := maker.extractClaim(token); ok {
		return claims.Username, nil
	}
	return "", fmt.Errorf("cannot extract username from jwt token")
}

func (maker *JWTMaker) GetRoleFromToken(token string) (string, error) {
	if claims, ok := maker.extractClaim(token); ok {
		return claims.Role, nil
	}
	return "", fmt.Errorf("cannot extract role form jwt token")
}

func (maker *JWTMaker) extractClaim(token string) (*UserPayload, bool) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidToken
		}
		return []byte(maker.secret), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &UserPayload{}, keyFunc)
	if err != nil {
		return nil, false
	}
	if claims, ok := jwtToken.Claims.(*UserPayload); ok {
		return claims, true
	}
	return nil, false
}
