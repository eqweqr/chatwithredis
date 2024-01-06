package service

import (
	"context"
	"errors"

	"github.com/eqweqr/chatwithredis/internal/dto"
	"github.com/eqweqr/chatwithredis/internal/repository"
	"github.com/eqweqr/chatwithredis/internal/token"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDosentExists  = errors.New("user not found")
)

type UserService interface {
	RegisterUser(context.Context, dto.RegisterRequest) (*dto.RegisterResponse, error)
	LoginUser(context.Context, dto.LoginRequest) (*dto.LoginResponse, error)
}

type UserServerImpl struct {
	tokenService   token.TokenMaker
	userRepository repository.UserRepository
}

func NewUserServer(tokenService token.TokenMaker, userStore repository.UserRepository) *UserServerImpl {
	return &UserServerImpl{
		tokenService:   tokenService,
		userRepository: userStore,
	}
}

// TODO: чделать через транзакцию rollback в случае если токен не сгенерировался.
func (server *UserServerImpl) RegisterUser(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	id, err := server.userRepository.AddUser(ctx, &req)
	if err != nil {
		return nil, err
	}

	token, _, err := server.tokenService.GenerateToken(req.Username, req.Role)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		Id:       id,
		Username: req.Username,
		Token:    token,
	}, nil
}

func (server *UserServerImpl) LoginUser(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := server.userRepository.GetUserByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, ErrUserDosentExists
	}
	if repository.CheckPassword(req.Password, user.EncryptedPassword) != nil {
		// поменять код ошибки.
		return nil, ErrUserDosentExists
	}

	token, _, err := server.tokenService.GenerateToken(req.Username, user.Role)
	if err != nil {
		return nil, ErrUserDosentExists
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil

}
