package httpserver

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const TokenExp = 30 * time.Minute

var SignMethod = jwt.SigningMethodHS256

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
	ID    int    `json:"id"`
}

func NewToken(login string, id int, secret []byte) (string, error) {
	token := jwt.NewWithClaims(SignMethod, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login: login,
		ID:    id,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func setToken(c echo.Context, token string) {
	resp := echo.NewResponse(c.Response().Writer, c.Echo())
	resp.Header().Add(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	c.SetResponse(resp)
}
