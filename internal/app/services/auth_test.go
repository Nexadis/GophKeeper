package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/Nexadis/GophKeeper/internal/models/users"
	mock_services "github.com/Nexadis/GophKeeper/mocks/intern/app/services"
	mock_models "github.com/Nexadis/GophKeeper/mocks/intern/models"
)

func TestNewAuth(t *testing.T) {
	type args struct {
		urepo UserRepo
		h     Hasher
		uf    UsersFactory
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
			if got := NewAuth(tt.args.urepo, tt.args.h, tt.args.uf); !reflect.DeepEqual(
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
		userRepo     *mock_services.MockUserRepo
		hasher       *mock_services.MockHasher
		userFactory  *mock_services.MockUsersFactory
		timeProvider *mock_models.MockTimeProvider
		cost         int
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
		want    users.User
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
				n := "username"
				p := "password"
				uf := users.New(f.timeProvider)
				u := uf.New(n, []byte(p))
				gomock.InOrder(
					f.hasher.EXPECT().Password(p).Return([]byte(p), nil),
					f.userFactory.EXPECT().New(n, []byte(p)).Return(u),
					f.userRepo.EXPECT().AddUser(nil, u).Return(nil),
					f.userRepo.EXPECT().GetUserByName(nil, n).Return(u, nil),
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
		{
			"AddUser problem",
			func(f *fields) {
				n := "username"
				p := "password"
				uf := users.New(f.timeProvider)
				u := uf.New(n, []byte(p))
				gomock.InOrder(
					f.hasher.EXPECT().Password(p).Return([]byte(p), nil),
					f.userFactory.EXPECT().New(n, []byte(p)).Return(u),
					f.userRepo.EXPECT().AddUser(nil, u).Return(fmt.Errorf("db is disconnected")),
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
			"GetUserByName problem",
			func(f *fields) {
				n := "username"
				p := "password"
				uf := users.New(f.timeProvider)
				u := uf.New(n, []byte(p))
				gomock.InOrder(
					f.hasher.EXPECT().Password(p).Return([]byte(p), nil),
					f.userFactory.EXPECT().New(n, []byte(p)).Return(u),
					f.userRepo.EXPECT().AddUser(nil, u).Return(nil),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_services.NewMockUserRepo(ctrl),
				mock_services.NewMockHasher(ctrl),
				mock_services.NewMockUsersFactory(ctrl),
				mock_models.NewMockTimeProvider(ctrl),
				defaultCost,
			}
			now := time.Now()
			if tt.prepare != nil {
				f.timeProvider.EXPECT().Now().Return(now).MinTimes(1)
				uf := users.New(f.timeProvider)
				if tt.wantErr != true {
					tt.want = uf.New(tt.args.username, []byte(tt.args.password))
				}
				tt.prepare(f)
			}
			a := &Auth{
				userRepo:    f.userRepo,
				hasher:      f.hasher,
				userFactory: f.userFactory,
				cost:        f.cost,
			}
			got, err := a.UserRegister(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.UserRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth.UserRegister() = %v, want %v", got, tt.want)
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
