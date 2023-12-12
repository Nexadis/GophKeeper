package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nexadis/GophKeeper/internal/models/users"
)

const defaultCost = bcrypt.DefaultCost

// UserRepo - позволяет работать с данными пользователей
type UserRepo interface {
	GetUserByName(ctx context.Context, username string) (*users.User, error)
	AddUser(ctx context.Context, u *users.User) error
	DeleteUser(ctx context.Context, username string) error
}

// Hasher - интерфейс для работы с паролем пользователя
type Hasher interface {
	Auth(ctx context.Context, u users.User, password string) error
	Password(password string) ([]byte, error)
}

// Auth - служба для авторизации пользователя
type Auth struct {
	userRepo UserRepo
	hasher   Hasher
	cost     int
}

// NewAuth - создаёт службу с заданными зависимостями
func NewAuth(urepo UserRepo, h Hasher) *Auth {
	cost := defaultCost
	return &Auth{
		urepo,
		h,
		cost,
	}
}

// UserRegister - регистрирует нового пользователя в системе
func (a *Auth) UserRegister(ctx context.Context, username, password string) (*users.User, error) {
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

// UserLogin - авторизует пользователя в системе
func (a *Auth) UserLogin(ctx context.Context, username, password string) (*users.User, error) {
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

// Password - преобразует пароль в форму для записи
func (h Hash) Password(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), h.cost)
}

// Auth - проверяет пароль на валидность для данного пользователя
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
