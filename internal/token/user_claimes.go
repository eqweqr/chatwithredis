// payload

package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ExpiredToken = errors.New("token has expired")
	InvalidToken = errors.New("token is invalid")
)

// instead standart claims
type UserPayload struct {
	IdToken   uuid.UUID `json:"id_token"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"` // время когда был выпущен токен.
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, role string, delay time.Duration) (*UserPayload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot create jwt token: %w", err)
	}

	claim := &UserPayload{
		IdToken:   id,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(delay),
	}

	return claim, nil
}

func (claim *UserPayload) Valid() error {
	if time.Now().After(claim.ExpiredAt) {
		return ExpiredToken
	}
	return nil
}
