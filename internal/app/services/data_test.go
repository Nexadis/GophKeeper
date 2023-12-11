package services

import (
	"context"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Nexadis/GophKeeper/internal/models/datas"
	mock_services "github.com/Nexadis/GophKeeper/mocks/intern/app/services"
)

func TestNewData(t *testing.T) {
	type args struct {
		drepo DataRepo
	}
	tests := []struct {
		name string
		args args
		want *Data
	}{
		{
			"Create Data Service",
			args{},
			&Data{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewData(tt.args.drepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_Add(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx   context.Context
		uid   int
		dlist []datas.Data
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Empty add list",
			nil,
			args{nil, 0, nil},
			false,
		},
		{
			"Add one invalid item without err",
			nil,
			args{nil, 1, []datas.Data{
				{
					Type:  datas.BankCardType,
					Value: "invalid",
				},
			}},
			true,
		},
		{
			"Add one valid item with err",
			func(f *fields) {
				f.dataRepo.EXPECT().Add(nil, gomock.Any()).Return(ErrAccessDenied)
			},
			args{nil, 1, []datas.Data{
				{
					Type:  datas.TextType,
					Value: "text",
				},
			}},
			true,
		},
		{
			"Add one valid item without err",
			func(f *fields) {
				f.dataRepo.EXPECT().Add(nil, gomock.Any()).Return(nil)
			},
			args{nil, 1, []datas.Data{
				{
					Type:  datas.TextType,
					Value: "text",
				},
			}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{mock_services.NewMockDataRepo(ctrl)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			d := &Data{
				dataRepo: f.dataRepo,
			}
			err := d.Add(tt.args.ctx, tt.args.uid, tt.args.dlist)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

	}
}

func TestData_Update(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx   context.Context
		uid   int
		dlist []datas.Data
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Empty update list",
			nil,
			args{nil, 0, nil},
			false,
		},
		{
			"Update invalid id item",
			nil,
			args{nil, 1, []datas.Data{
				{
					Type:  datas.BankCardType,
					Value: "invalid",
				},
			}},
			true,
		},
		{
			"Update invalid item value",
			nil,
			args{nil, 1, []datas.Data{
				{
					ID:    1,
					Type:  datas.BankCardType,
					Value: "invalid",
				},
			}},
			true,
		},
		{
			"Update error in dataRepo",
			func(f *fields) {
				f.dataRepo.EXPECT().Update(nil, gomock.Any()).Return(ErrDataNotFound)
			},
			args{nil, 1, []datas.Data{
				{
					ID:    1,
					Type:  datas.TextType,
					Value: "text",
				},
			}},
			true,
		},
		{
			"Update in dataRepo",
			func(f *fields) {
				f.dataRepo.EXPECT().Update(nil, gomock.Any()).Return(nil)
			},
			args{nil, 1, []datas.Data{
				{
					ID:    1,
					Type:  datas.TextType,
					Value: "text",
				},
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			if err := ds.Update(tt.args.ctx, tt.args.uid, tt.args.dlist); (err != nil) != tt.wantErr {
				t.Errorf("Data.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_GetByID(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx context.Context
		uid int
		id  int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *datas.Data
		wantErr bool
	}{
		{
			"Not found id",
			func(f *fields) {
				f.dataRepo.EXPECT().GetByID(nil, 1).Return(nil, ErrDataNotFound)
			},
			args{nil, 2, 1},
			nil,
			true,
		},
		{
			"Don't match uid and data.UserID",
			func(f *fields) {
				d, _ := datas.NewData(datas.TextType, "text")
				d.UserID = 1
				f.dataRepo.EXPECT().GetByID(nil, 1).Return(d, nil)
			},
			args{nil, 2, 1},
			nil,
			true,
		},
		{
			"Data found",
			func(f *fields) {
				d, _ := datas.NewData(datas.TextType, "text")
				d.UserID = 2
				f.dataRepo.EXPECT().GetByID(nil, 1).Return(d, nil)
			},
			args{nil, 2, 1},
			&datas.Data{
				Type:  datas.TextType,
				Value: "text",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := Data{
				dataRepo: f.dataRepo,
			}
			got, err := ds.GetByID(tt.args.ctx, tt.args.uid, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !(got.Type == tt.want.Type && got.Value == tt.want.Value) {
				t.Errorf("Data.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_GetByUser(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx context.Context
		uid int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    []datas.Data
		wantErr bool
	}{
		{
			"Can't fetch data",
			func(f *fields) {
				f.dataRepo.EXPECT().GetByUser(nil, 1).Return(nil, ErrInvalidUserID)
			},
			args{nil, 1},
			nil,
			true,
		},
		{
			"Fetched user",
			func(f *fields) {
				f.dataRepo.EXPECT().GetByUser(nil, 1).Return(nil, nil)
			},
			args{nil, 1},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			ds := Data{
				dataRepo: f.dataRepo,
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			got, err := ds.GetByUser(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.GetByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_DeleteByID(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx context.Context
		uid int
		id  []int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Empty list",
			nil,
			args{},
			false,
		},
		{
			"Delete normal",
			func(f *fields) {
				f.dataRepo.EXPECT().DeleteByIDs(nil, 1, gomock.Any()).Return(nil)
			},
			args{nil, 1, []int{2}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			if err := ds.DeleteByID(tt.args.ctx, tt.args.uid, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Data.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_Health(t *testing.T) {
	type fields struct {
		dataRepo DataRepo
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &Data{
				dataRepo: tt.fields.dataRepo,
			}
			if err := ds.Health(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Data.Health() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
