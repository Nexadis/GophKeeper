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
	Add(ctx context.Context, datas []datas.Data) error
	GetByID(ctx context.Context, id int) (*datas.Data, error)
	GetByUser(ctx context.Context, uid int) ([]datas.Data, error)
	Update(ctx context.Context, data []datas.Data) error
	DeleteByIDs(ctx context.Context, uid int, ids []int) error
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

func (ds *Data) Add(ctx context.Context, uid int, dlist []datas.Data) error {
	for i, d := range dlist {
		err := dlist[i].SetValue(d.Value)
		if err != nil {
			return fmt.Errorf("can't add data '%s' for uid=%d : %w'", d.Value, uid, err)
		}
		d.SetValue(d.Value)
		dlist[i].UserID = uid
		dlist[i].CreatedAt = dlist[i].EditedAt

	}
	err := ds.dataRepo.Add(ctx, dlist)
	if err != nil {
		return fmt.Errorf("can't add data for uid=%d : %w", uid, err)
	}
	return nil
}

func (ds *Data) Update(ctx context.Context, uid int, dlist []datas.Data) error {
	for i, d := range dlist {
		if d.ID == 0 {
			return ErrInvalidDataID
		}
		err := d.SetValue(d.Value)
		if err != nil {
			return fmt.Errorf("can't update data '%s' with id=%d : %w'", d.Value, d.ID, err)
		}
		d.UserID = uid
		dlist[i] = d

	}
	err := ds.dataRepo.Update(ctx, dlist)
	if err != nil {
		return fmt.Errorf("can't update data : %w'", err)
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

func (ds Data) GetByUser(ctx context.Context, uid int) ([]datas.Data, error) {
	datas, err := ds.dataRepo.GetByUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("can't get data for uid=%d : %w", uid, err)
	}
	return datas, nil
}

func (ds *Data) DeleteByID(ctx context.Context, uid int, id []int) error {
	return ds.dataRepo.DeleteByIDs(ctx, uid, id)
}

func (ds *Data) Health(ctx context.Context) error {
	err := ds.dataRepo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Data repo is unavailable: %w", err)
	}
	return nil
}
