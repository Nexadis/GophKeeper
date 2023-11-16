package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nexadis/GophKeeper/internal/users"
)

var (
	ErrDataNotFound    = errors.New("data not found")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidDataID   = errors.New("invalid data id")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrAccessDenied    = errors.New("access denied")
)

type IData interface {
	ID() int
	SetID(id int)
	UserID() int
	SetUserID(id int)
	Description() string
	SetDescription(desc string)
	CreatedAt() time.Time
	EditedAt() time.Time
	SetValue(value string) error
	Value() string
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (users.User, error)
	GetUserByName(ctx context.Context, username string) (users.User, error)
	AddUser(ctx context.Context, u users.User) error
	DeleteUser(ctx context.Context, username string) error
}

type (
	Auth interface {
		Auth(ctx context.Context, u users.User, password string) error
		Hash(password string) ([]byte, error)
	}
)

type DataRepo interface {
	Add(ctx context.Context, data IData) error
	GetByID(ctx context.Context, id int) (IData, error)
	GetByUser(ctx context.Context, u users.User) ([]IData, error)
	Update(ctx context.Context, data IData) error
	DeleteByID(ctx context.Context, id int) error
}

type app struct {
	userRepo UserRepo
	dataRepo DataRepo
	auth     Auth
}

func New(ds DataRepo, us UserRepo, auth Auth) app {
	return app{
		us,
		ds,
		auth,
	}
}

func (a *app) Add(ctx context.Context, u users.User, data IData) error {
	data.SetUserID(u.ID())
	err := a.dataRepo.Add(ctx, data)
	if err != nil {
		return fmt.Errorf("can't add data '%s' for %s : %w", data.Value(), u.Username(), err)
	}
	return nil
}

func (a *app) Update(ctx context.Context, u users.User, data IData) error {
	if data.ID() == 0 {
		return ErrInvalidDataID
	}
	d, err := a.dataRepo.GetByID(ctx, data.ID())
	if err != nil {
		return fmt.Errorf("can't update data '%s' : %w", data.Value(), err)
	}
	if d.UserID() != u.ID() {
		return fmt.Errorf("you can't update this data '%s' : %w", data.Value(), ErrAccessDenied)
	}
	err = a.dataRepo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("can't update data '%s' : %w", data.Value(), err)
	}

	return nil
}

func (a *app) GetByID(ctx context.Context, u users.User, id int) (IData, error) {
	d, err := a.dataRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get data with id %d : %w", id, err)
	}
	if d.UserID() != u.ID() {
		return nil, fmt.Errorf("you aren't owner of data with id %d : %w", id, ErrAccessDenied)
	}
	return d, nil
}

func (a *app) GetByUser(ctx context.Context, u users.User) ([]IData, error) {
	datas, err := a.dataRepo.GetByUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("can't get data for user %s : %w", u.Username(), err)
	}
	return datas, nil
}

func (a *app) DeleteByID(ctx context.Context, u users.User, id int) error {
	d, err := a.dataRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete data with id %d : %w", id, err)
	}
	if d.UserID() != u.ID() {
		return fmt.Errorf("can't delete data with id %d : %w", id, ErrAccessDenied)
	}
	a.dataRepo.DeleteByID(ctx, id)
	return nil
}

func (a *app) UserRegister(ctx context.Context, username, password string) (users.User, error) {
	if !a.validUsername(username) {
		return nil, ErrInvalidUsername
	}
	if !a.validPassword(password) {
		return nil, ErrInvalidUsername
	}
	hash, err := a.auth.Hash(password)
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

func (a *app) UserLogin(ctx context.Context, username, password string) (users.User, error) {
	u, err := a.userRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("can't login user %s : %w", username, err)
	}
	err = a.auth.Auth(ctx, u, password)
	if err != nil {
		return nil, fmt.Errorf("can't login user %s : %w", username, err)
	}
	return u, err
}

func (a app) validPassword(password string) bool {
	return true
}

func (a app) validUsername(username string) bool {
	return true
}
