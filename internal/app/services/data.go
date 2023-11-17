package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nexadis/GophKeeper/internal/models/users"
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

type DataRepo interface {
	Add(ctx context.Context, data IData) error
	GetByID(ctx context.Context, id int) (IData, error)
	GetByUser(ctx context.Context, u users.User) ([]IData, error)
	Update(ctx context.Context, data IData) error
	DeleteByID(ctx context.Context, id int) error
}

type Data struct {
	dataRepo DataRepo
}

func NewData(drepo DataRepo) Data {
	return Data{
		drepo,
	}
}

func (ds *Data) Add(ctx context.Context, u users.User, data IData) error {
	data.SetUserID(u.ID())
	err := ds.dataRepo.Add(ctx, data)
	if err != nil {
		return fmt.Errorf("can't add data '%s' for %s : %w", data.Value(), u.Username(), err)
	}
	return nil
}

func (ds *Data) Update(ctx context.Context, u users.User, data IData) error {
	if data.ID() == 0 {
		return ErrInvalidDataID
	}
	d, err := ds.dataRepo.GetByID(ctx, data.ID())
	if err != nil {
		return fmt.Errorf("can't update data '%s' : %w", data.Value(), err)
	}
	if d.UserID() != u.ID() {
		return fmt.Errorf("you can't update this data '%s' : %w", data.Value(), ErrAccessDenied)
	}
	err = ds.dataRepo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("can't update data '%s' : %w", data.Value(), err)
	}

	return nil
}

func (ds Data) GetByID(ctx context.Context, u users.User, id int) (IData, error) {
	d, err := ds.dataRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get data with id %d : %w", id, err)
	}
	if d.UserID() != u.ID() {
		return nil, fmt.Errorf("you aren't owner of data with id %d : %w", id, ErrAccessDenied)
	}
	return d, nil
}

func (ds Data) GetByUser(ctx context.Context, u users.User) ([]IData, error) {
	datas, err := ds.dataRepo.GetByUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("can't get data for user %s : %w", u.Username(), err)
	}
	return datas, nil
}

func (ds *Data) DeleteByID(ctx context.Context, u users.User, id int) error {
	d, err := ds.dataRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete data with id %d : %w", id, err)
	}
	if d.UserID() != u.ID() {
		return fmt.Errorf("can't delete data with id %d : %w", id, ErrAccessDenied)
	}
	ds.dataRepo.DeleteByID(ctx, id)
	return nil
}
