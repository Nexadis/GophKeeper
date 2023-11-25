package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nexadis/GophKeeper/internal/models/users"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (users.User, error)
	GetUserByName(ctx context.Context, username string) (users.User, error)
	AddUser(ctx context.Context, u users.User) error
	DeleteUser(ctx context.Context, username string) error
}

type Auther interface {
	Auth(ctx context.Context, u users.User, password string) error
	Hash(password string) ([]byte, error)
}

type Auth struct {
	userRepo UserRepo
	cost     int
}

func NewAuth(urepo UserRepo, cost int) Auth {
	return Auth{
		urepo,
		cost,
	}
}

func (a *Auth) UserRegister(
	ctx context.Context,
	username, password string,
) (users.User, error) {
	if !a.validUsername(username) {
		return nil, ErrInvalidUsername
	}
	if !a.validPassword(password) {
		return nil, ErrInvalidUsername
	}
	hash, err := a.hash(password)
	if err != nil {
		return nil, fmt.Errorf("can't register user: %w", err)
	}
	u := users.New(username, hash)
	err = a.userRepo.AddUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("can't register user: %w", err)
	}
	return a.UserLogin(ctx, username, password)
}

func (a *Auth) UserLogin(
	ctx context.Context,
	username, password string,
) (users.User, error) {
	u, err := a.userRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("can't login user %s : %w", username, err)
	}
	err = a.auth(ctx, u, password)
	if err != nil {
		return nil, fmt.Errorf("can't login user %s : %w", username, err)
	}
	return u, err
}

func (a Auth) validPassword(password string) bool {
	return true
}

func (a Auth) validUsername(username string) bool {
	return true
}

func (a Auth) hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), a.cost)
}

func (a Auth) auth(ctx context.Context, u users.User, password string) error {
	return bcrypt.CompareHashAndPassword(u.Hash(), []byte(password))
}
