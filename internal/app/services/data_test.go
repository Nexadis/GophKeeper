package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Nexadis/GophKeeper/internal/models/datas"
	"github.com/Nexadis/GophKeeper/internal/models/users"
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
	}{{
		"Data service",
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
		ctx  context.Context
		u    users.User
		data datas.IData
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		wantErr bool
	}{
		{
			"Valid user and data",
			func(f *fields) {
				f.dataRepo.EXPECT().Add(nil, gomock.Any()).Return(nil)
			},
			args{
				nil,
				users.New(nil).New("username", []byte("password")),
				datas.NewCredentials("save_login", "save_password"),
			},
			false,
		},
		{
			"Error Add user",
			func(f *fields) {
				f.dataRepo.EXPECT().Add(nil, gomock.Any()).Return(fmt.Errorf("db is disconnected"))
			},
			args{
				nil,
				users.New(nil).New("username", []byte("password")),
				datas.NewCredentials("save_login", "save_password"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			if err := ds.Add(tt.args.ctx, tt.args.u, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Data.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_Update(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx  context.Context
		u    users.User
		data datas.IData
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		wantErr bool
	}{
		{
			"Error data id",
			nil,
			func() args {
				u := users.New(nil).New("username", []byte("password"))
				d := datas.NewCredentials("save_login", "save_password")
				return args{
					nil,
					u,
					d,
				}
			},
			true,
		},
		{
			"Error Get data with id",
			func(f *fields) {
				gomock.InOrder(
					f.dataRepo.EXPECT().
						GetByID(nil, 123).
						Return(nil, fmt.Errorf("data is not exist")),
				)
			},
			func() args {
				u := users.New(nil).New("username", []byte("password"))
				d := datas.NewCredentials("save_login", "save_password")
				d.SetID(123)
				return args{
					nil,
					u,
					d,
				}
			},
			true,
		},
		{
			"Error AccessDenied",
			func(f *fields) {
				gomock.InOrder(
					f.dataRepo.EXPECT().
						GetByID(nil, 123).
						Return(datas.NewCredentials("u", "p"), nil),
				)
			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(1)
				d := datas.NewCredentials("save_login", "save_password")
				d.SetID(123)
				d.SetUserID(2)
				return args{
					nil,
					u,
					d,
				}
			},
			true,
		},
		{
			"Error Update",
			func(f *fields) {
				d := datas.NewCredentials("u", "p")
				d.SetID(123)
				d.SetUserID(1)
				gomock.InOrder(
					f.dataRepo.EXPECT().
						GetByID(nil, 123).
						Return(d, nil),
					f.dataRepo.EXPECT().
						Update(nil, gomock.Any()).
						Return(fmt.Errorf("db is disconnected")),
				)
			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(1)
				d := datas.NewCredentials("save_login", "save_password")
				d.SetID(123)
				d.SetUserID(1)
				return args{
					nil,
					u,
					d,
				}
			},
			true,
		},
		{
			"Good way",
			func(f *fields) {
				d := datas.NewCredentials("u", "p")
				d.SetID(123)
				d.SetUserID(1)
				gomock.InOrder(
					f.dataRepo.EXPECT().
						GetByID(nil, 123).
						Return(d, nil),
					f.dataRepo.EXPECT().
						Update(nil, gomock.Any()).
						Return(nil),
				)
			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(1)
				d := datas.NewCredentials("save_login", "save_password")
				d.SetID(123)
				d.SetUserID(1)
				return args{
					nil,
					u,
					d,
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			args := tt.setargs()
			if err := ds.Update(args.ctx, args.u, args.data); (err != nil) != tt.wantErr {
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
		u   users.User
		id  int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		want    datas.IData
		wantErr bool
	}{
		{
			"Error with db",
			func(f *fields) {
				f.dataRepo.EXPECT().GetByID(nil, 1).Return(nil, fmt.Errorf("db is disconnected"))

			},
			func() args {
				return args{
					nil,
					nil,
					1,
				}
			},
			nil,
			true,
		},
		{
			"Error AccessDenied",
			func(f *fields) {
				d := datas.NewText("text data")
				d.SetUserID(1)
				f.dataRepo.EXPECT().GetByID(nil, 3).Return(d, nil)

			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(2)
				return args{
					nil,
					u,
					3,
				}
			},
			nil,
			true,
		},
		{
			"Valid data",
			func(f *fields) {
				d := datas.NewText("text data")
				d.SetUserID(1)
				f.dataRepo.EXPECT().GetByID(nil, 3).Return(d, nil)

			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(1)
				return args{
					nil,
					u,
					3,
				}
			},
			datas.NewText("text data"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			args := tt.setargs()
			got, err := ds.GetByID(args.ctx, args.u, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !reflect.DeepEqual(got.Value(), tt.want.Value()) {
					t.Errorf("Data.GetByID() = %v, want %v", got, tt.want)
				}
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
		u   users.User
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		want    []datas.IData
		wantErr bool
	}{
		{
			"Error with db",
			func(f *fields) {
				f.dataRepo.EXPECT().
					GetByUser(nil, gomock.Any()).
					Return(nil, fmt.Errorf("user not found"))
			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				return args{
					nil,
					u,
				}
			},
			nil,
			true,
		},
		{
			"Valid User",
			func(f *fields) {
				ds := []datas.IData{
					datas.NewText("some"),
				}
				f.dataRepo.EXPECT().
					GetByUser(nil, gomock.Any()).
					Return(ds, nil)
			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				return args{
					nil,
					u,
				}
			},
			[]datas.IData{
				datas.NewText("some"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			args := tt.setargs()
			got, err := ds.GetByUser(args.ctx, args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if len(tt.want) != len(got) {
					t.Errorf("Data.GetByUser() = len %v, want len %v", len(got), len(tt.want))
					return

				}
				for i, v := range tt.want {
					if !reflect.DeepEqual(got[i].Value(), v.Value()) {
						t.Errorf(
							"Data.GetByUser() [%d] = %v, want %v",
							i,
							got[i].Value(),
							v.Value(),
						)
					}

				}

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
		u   users.User
		id  int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		wantErr bool
	}{
		{
			"Error with db",
			func(f *fields) {
				f.dataRepo.EXPECT().GetByID(nil, 1).Return(nil, fmt.Errorf("db is disconnected"))

			},
			func() args {
				return args{
					nil,
					nil,
					1,
				}
			},
			true,
		},
		{
			"Error Access Denied",
			func(f *fields) {
				d := datas.NewText("text data")
				d.SetUserID(1)
				f.dataRepo.EXPECT().GetByID(nil, 3).Return(d, nil)

			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(2)
				return args{
					nil,
					u,
					3,
				}
			},
			true,
		},
		{
			"Valid data",
			func(f *fields) {
				d := datas.NewText("text data")
				d.SetUserID(1)
				f.dataRepo.EXPECT().GetByID(nil, 3).Return(d, nil)
				f.dataRepo.EXPECT().DeleteByID(nil, 3).Return(nil)

			},
			func() args {
				u := users.New(nil).New("u", []byte("p"))
				u.SetID(1)
				return args{
					nil,
					u,
					3,
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			args := tt.setargs()
			if err := ds.DeleteByID(args.ctx, args.u, args.id); (err != nil) != tt.wantErr {
				t.Errorf("Data.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_Health(t *testing.T) {
	type fields struct {
		dataRepo *mock_services.MockDataRepo
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		wantErr bool
	}{
		{
			"Error with db",
			func(f *fields) {
				f.dataRepo.EXPECT().Ping(nil).Return(fmt.Errorf("db is disconnected"))
			},
			func() args {
				return args{
					nil,
				}
			},
			true,
		},
		{
			"Ping is work",
			func(f *fields) {
				f.dataRepo.EXPECT().Ping(nil).Return(nil)
			},
			func() args {
				return args{
					nil,
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockDataRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			ds := &Data{
				dataRepo: f.dataRepo,
			}
			args := tt.setargs()
			if err := ds.Health(args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Data.Health() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
