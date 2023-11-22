package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Nexadis/GophKeeper/internal/models/users"
	mock_services "github.com/Nexadis/GophKeeper/mocks/intern/app/services"
)

func TestNewAuth(t *testing.T) {
	type args struct {
		urepo UserRepo
		h     Hasher
	}
	tests := []struct {
		name string
		args args
		want *Auth
	}{{
		"Empty Auth",
		args{},
		&Auth{
			cost: defaultCost,
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuth(tt.args.urepo, tt.args.h); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_UserRegister(t *testing.T) {
	type fields struct {
		userRepo *mock_services.MockUserRepo
		hasher   *mock_services.MockHasher
		cost     int
	}
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *users.User
		wantErr bool
	}{
		{
			"Empty username",
			nil,
			args{
				nil,
				"",
				"password",
			},
			nil,
			true,
		},
		{
			"Empty password",
			nil,
			args{
				nil,
				"username",
				"",
			},
			nil,
			true,
		},
		{
			"Normal user",
			func(f *fields) {
				p := "password"
				gomock.InOrder(
					f.hasher.EXPECT().Password(p).Return([]byte(p), nil),
					f.userRepo.EXPECT().AddUser(nil, gomock.Any()).Return(nil),
				)

			},
			args{
				nil,
				"username",
				"password",
			},
			nil,
			false,
		},
		{
			"AddUser problem",
			func(f *fields) {
				p := "password"
				gomock.InOrder(
					f.hasher.EXPECT().Password(p).Return([]byte(p), nil),
					f.userRepo.EXPECT().
						AddUser(nil, gomock.Any()).
						Return(fmt.Errorf("db is disconnected")),
				)

			},
			args{
				nil,
				"username",
				"password",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockUserRepo(ctrl),
				mock_services.NewMockHasher(ctrl),
				defaultCost,
			}
			if tt.prepare != nil {
				if tt.wantErr != true {
					u := users.New(tt.args.username, []byte(tt.args.password))
					tt.want = &u
				}
				tt.prepare(f)
			}
			a := &Auth{
				userRepo: f.userRepo,
				hasher:   f.hasher,
				cost:     f.cost,
			}
			got, err := a.UserRegister(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.UserRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !(got.Username == tt.want.Username && reflect.DeepEqual(got.Hash, tt.want.Hash) &&
					got.ID == tt.want.ID) {
					t.Errorf("Auth.UserRegister() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
func TestAuth_UserLogin(t *testing.T) {
	type fields struct {
		userRepo *mock_services.MockUserRepo
		hasher   *mock_services.MockHasher
		cost     int
	}
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *users.User
		wantErr bool
	}{
		{
			"GetUserByName problem",
			func(f *fields) {
				n := "username"
				gomock.InOrder(
					f.userRepo.EXPECT().
						GetUserByName(nil, n).
						Return(nil, fmt.Errorf("db is disconnected")),
				)

			},
			args{
				nil,
				"username",
				"password",
			},
			nil,
			true,
		},
		{
			"Auth problem",
			func(f *fields) {
				n := "username"
				p := "password"
				u := users.New(n, []byte(p))
				gomock.InOrder(
					f.userRepo.EXPECT().
						GetUserByName(nil, n).
						Return(&u, nil),
					f.hasher.EXPECT().Auth(nil, u, p).Return(ErrAccessDenied),
				)

			},
			args{
				nil,
				"username",
				"password",
			},
			nil,
			true,
		},
		{
			"Good Login",
			func(f *fields) {
				n := "username"
				p := "password"
				u := users.New(n, []byte(p))
				gomock.InOrder(
					f.userRepo.EXPECT().
						GetUserByName(nil, n).
						Return(&u, nil),
					f.hasher.EXPECT().Auth(nil, u, p).Return(nil),
				)

			},
			args{
				nil,
				"username",
				"password",
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockUserRepo(ctrl),
				mock_services.NewMockHasher(ctrl),
				defaultCost,
			}
			if tt.prepare != nil {
				if tt.wantErr != true {
					u := users.New(tt.args.username, []byte(tt.args.password))
					tt.want = &u
				}
				tt.prepare(f)
			}
			a := &Auth{
				userRepo: f.userRepo,
				hasher:   f.hasher,
				cost:     f.cost,
			}
			got, err := a.UserLogin(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.UserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !(got.Username == tt.want.Username && reflect.DeepEqual(got.Hash, tt.want.Hash) &&
					got.ID == tt.want.ID) {
					t.Errorf("Auth.UserLogin() = %v, want %v", got, tt.want)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
			}
		})
	}
}

func TestNewHash(t *testing.T) {
	tests := []struct {
		name string
		want *Hash
	}{
		{
			"Just Hash",
			&Hash{
				defaultCost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
