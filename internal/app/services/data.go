package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

var (
	ErrDataNotFound    = errors.New("data not found")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidDataID   = errors.New("invalid data id")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrAccessDenied    = errors.New("access denied")
)

type DataRepo interface {
	Add(ctx context.Context, data *datas.Data) error
	GetByID(ctx context.Context, id int) (*datas.Data, error)
	GetByUser(ctx context.Context, uid int) ([]*datas.Data, error)
	Update(ctx context.Context, data *datas.Data) error
	DeleteByID(ctx context.Context, id int) error
	Ping(ctx context.Context) error
}

type Data struct {
	dataRepo DataRepo
}

func NewData(drepo DataRepo) *Data {
	return &Data{
		drepo,
	}
}

func (ds *Data) Add(ctx context.Context, uid int, data *datas.Data) error {
	d, err := datas.NewData(data.Type, data.Value)
	if err != nil {
		return fmt.Errorf("can't add data '%s' for uid=%d : %w'", data.Value, uid, err)
	}
	d.UserID = uid
	err = ds.dataRepo.Add(ctx, d)
	if err != nil {
		return fmt.Errorf("can't add data '%s' for uid=%d : %w", data.Value, uid, err)
	}
	*data = *d
	return nil
}

func (ds *Data) Update(ctx context.Context, uid int, data *datas.Data) error {
	if data.ID == 0 {
		return ErrInvalidDataID
	}
	err := data.SetValue(data.Value)
	if err != nil {
		return fmt.Errorf("can't update data '%s' with id=%d : %w'", data.Value, data.ID, err)
	}
	d, err := ds.dataRepo.GetByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("can't update data '%s' with id=%d : %w'", data.Value, data.ID, err)
	}
	if d.UserID != uid {
		return fmt.Errorf(
			"uid=%d can't update this data '%s' : %w",
			uid,
			data.Value,
			ErrAccessDenied,
		)
	}
	err = ds.dataRepo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("can't update data '%s' with id=%d : %w'", data.Value, data.ID, err)
	}

	return nil
}

func (ds Data) GetByID(ctx context.Context, uid int, id int) (*datas.Data, error) {
	d, err := ds.dataRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get data with id=%d : %w", id, err)
	}
	if d.UserID != uid {
		return nil, fmt.Errorf(
			"uid=%d isn't owner of data with id=%d : %w",
			uid,
			id,
			ErrAccessDenied,
		)
	}
	return d, nil
}

func (ds Data) GetByUser(ctx context.Context, uid int) ([]*datas.Data, error) {
	datas, err := ds.dataRepo.GetByUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("can't get data for uid=%d : %w", uid, err)
	}
	return datas, nil
}

func (ds *Data) DeleteByID(ctx context.Context, uid, id int) error {
	d, err := ds.dataRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete data with id %d : %w", id, err)
	}
	if d.UserID != uid {
		return fmt.Errorf("uid=%d can't delete data with id %d : %w", uid, id, ErrAccessDenied)
	}
	return ds.dataRepo.DeleteByID(ctx, id)
}

func (ds *Data) Health(ctx context.Context) error {
	err := ds.dataRepo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Data repo is unavailable: %w", err)
	}
	return nil
}
