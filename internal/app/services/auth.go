package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nexadis/GophKeeper/internal/models/users"
)

const defaultCost = bcrypt.DefaultCost

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (*users.User, error)
	GetUserByName(ctx context.Context, username string) (*users.User, error)
	AddUser(ctx context.Context, u *users.User) error
	DeleteUser(ctx context.Context, username string) error
}

type Hasher interface {
	Auth(ctx context.Context, u users.User, password string) error
	Password(password string) ([]byte, error)
}

type Auth struct {
	userRepo UserRepo
	hasher   Hasher
	cost     int
}

func NewAuth(urepo UserRepo, h Hasher) *Auth {
	cost := defaultCost
	return &Auth{
		urepo,
		h,
		cost,
	}
}

func (a *Auth) UserRegister(
	ctx context.Context,
	username, password string,
) (*users.User, error) {
	if !a.validUsername(username) {
		return nil, ErrInvalidUsername
	}
	if !a.validPassword(password) {
		return nil, ErrInvalidPassword
	}
	hash, err := a.hasher.Password(password)
	if err != nil {
		return nil, fmt.Errorf("can't register user %s : %w", username, err)
	}
	u := users.New(username, hash)
	err = a.userRepo.AddUser(ctx, &u)
	if err != nil {
		return nil, fmt.Errorf("can't register user %s : %w", username, err)
	}
	return &u, nil
}

func (a *Auth) UserLogin(
	ctx context.Context,
	username, password string,
) (*users.User, error) {
	u, err := a.userRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("can't login user %s : %w", username, err)
	}
	err = a.hasher.Auth(ctx, *u, password)
	if err != nil {
		return nil, fmt.Errorf("can't login user %s : %w", username, err)
	}
	return u, err
}

func (a Auth) validPassword(password string) bool {
	if password == "" {
		return false
	}
	return true
}

func (a Auth) validUsername(username string) bool {
	if username == "" {
		return false
	}
	return true
}

func (h Hash) Password(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), h.cost)
}

func (h Hash) Auth(ctx context.Context, u users.User, password string) error {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password))
}

type Hash struct {
	cost int
}

func NewHash() *Hash {
	return &Hash{
		defaultCost,
	}
}
