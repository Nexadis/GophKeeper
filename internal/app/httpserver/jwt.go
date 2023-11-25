package httpserver

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/Nexadis/GophKeeper/internal/logger"
)

const TokenExp = 30 * time.Minute

var SignMethod = jwt.SigningMethodHS256

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
	UID   int    `json:"uid"`
}

func NewToken(login string, uid int, secret []byte) (string, error) {
	token := jwt.NewWithClaims(SignMethod, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login: login,
		UID:   uid,
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

func GetUID(c echo.Context) (int, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("can't map jwt")

	}
	logger.Debug(fmt.Sprintf("Got claims: %v", claims))

	return claims.UID, nil
}
