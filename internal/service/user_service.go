package service

import (
	"context"

	"github.com/eqweqr/chatwithredis/internal/dto"
	"github.com/eqweqr/chatwithredis/internal/model"
)

type UserService interface {
	RegisterUser(context.Context, dto.RegisterRequest) (*model.User, error)
}
