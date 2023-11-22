package httpserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Nexadis/GophKeeper/internal/database"
	"github.com/Nexadis/GophKeeper/internal/logger"
)

func (hs *Server) Register(c echo.Context) error {
	u := &User{}
	err := c.Bind(u)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	user, err := hs.authService.UserRegister(c.Request().Context(), u.Login, u.Password)
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, database.ErrUserExist):
			return c.String(http.StatusConflict, database.ErrUserExist.Error())

		}
		return c.NoContent(http.StatusInternalServerError)
	}
	token, err := NewToken(user.Username, user.ID, hs.config.JWTSecret)
	if err != nil {
		logger.Error(fmt.Errorf("Can't generate token for user %s: %w", u.Login, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	setToken(c, token)

	return c.NoContent(http.StatusOK)
}
func (hs *Server) Login(c echo.Context) error {
	u := &User{}
	err := c.Bind(u)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	user, err := hs.authService.UserLogin(c.Request().Context(), u.Login, u.Password)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusForbidden, err)
	}
	token, err := NewToken(user.Username, user.ID, hs.config.JWTSecret)
	if err != nil {
		logger.Error(fmt.Errorf("Can't generate token for user %s: %w", u.Login, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	setToken(c, token)

	return c.NoContent(http.StatusOK)
}

func (hs *Server) GetData(c echo.Context) error {
	return nil
}
func (hs *Server) PostData(c echo.Context) error {
	return nil
}
func (hs *Server) DeleteData(c echo.Context) error {
	return nil
}
func (hs *Server) UpdateData(c echo.Context) error {
	return nil
}

func (hs *Server) GetUserData(c echo.Context) error {
	return nil
}
func (hs *Server) PostUserData(c echo.Context) error {
	return nil
}
func (hs *Server) DeleteUserData(c echo.Context) error {
	return nil
}
func (hs *Server) UpdateUserData(c echo.Context) error {
	return nil
}

func (hs *Server) Ping(c echo.Context) error {
	err := hs.dataService.Health(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return nil
}
