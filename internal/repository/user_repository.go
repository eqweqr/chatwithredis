package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eqweqr/chatwithredis/internal/db"
	"github.com/eqweqr/chatwithredis/internal/dto"
	"github.com/eqweqr/chatwithredis/internal/model"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *dto.RegisterRequest) (int, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	DeleteUserByID(ctx context.Context, id int) error
}

type UserStore struct {
	db db.UserDB
}

func NewUserStore(db db.UserDB) *UserStore {
	return &UserStore{db}
}

func (store *UserStore) DeleteUserByID(ctx context.Context, id int) error {
	query := `delete from users where id=$1`
	_, err := store.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}
	return nil
}

func (store *UserStore) AddUser(ctx context.Context, userDTO *dto.RegisterRequest) (int, error) {
	query := `INSERT INTO users(username, encrypted_password, role) values($1, $2, $3) returning id`
	encryptedPass, err := EncryptPassword(userDTO.Password)
	if err != nil {
		return -1, fmt.Errorf("cannot generate password: %w", err)
	}
	result, err := store.db.QueryRowContext(ctx, query, userDTO.Username, encryptedPass, userDTO.Role)
	if err != nil {
		return -1, fmt.Errorf("cannot add user to db: %w", err)
	}
	var id int
	if err := result.Scan(&id); err != nil {
		return -1, fmt.Errorf("cannot get id: %w", err)
	}

	return id, nil
}

func (store *UserStore) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	query := `select * from users where id=$1`
	result, err := store.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("cannot extract user: %w", err)
	}
	var user *model.User
	if err := result.Scan(user.ID, user.Username, user.EncryptedPassword, user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot scan user with id(%d): %w", id, err)
	}
	return user, nil
}

func (store *UserStore) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `select * from users where username=$1`
	result, err := store.db.QueryContext(ctx, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot extract user: %w", err)
	}
	var user *model.User
	if err := result.Scan(user.ID, user.Username, user.EncryptedPassword, user.Role); err != nil {
		return nil, fmt.Errorf("cannot scan user with id(%s): %w", username, err)
	}
	return user, nil
}
func (store *UserStore) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	query := `select * from users`
	result, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("cannot extract users from db: %w", err)
	}
	var users []*model.User
	for result.Next() {
		var user *model.User
		err := result.Scan(user.ID, user.Username, user.EncryptedPassword, user.Role)
		if err != nil {
			return nil, fmt.Errorf("cannot ")
		}
		users = append(users, user)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("cannot in result set: %w", err)
	}
	return users, nil
}
