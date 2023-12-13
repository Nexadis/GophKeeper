package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		username string
		hash     []byte
	}
	tests := []struct {
		name string
		args args
		want User
	}{{
		"Just User",
		args{
			"Login",
			[]byte("hash"),
		},
		User{
			Username: "Login",
			Hash:     []byte("hash"),
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.username, tt.args.hash); !assert.Equal(
				t,
				got.Username,
				tt.want.Username,
			) &&
				!assert.Equal(t, got.Hash, tt.want.Hash) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
