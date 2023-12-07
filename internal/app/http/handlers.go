package http

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/Nexadis/GophKeeper/internal/app/services"
	"github.com/Nexadis/GophKeeper/internal/database"
	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

func (hs *Server) Download(c echo.Context) error {
	filename := c.Param("file")
	logger.Infof("Download %s", filename)
	d, err := os.ReadDir(hs.config.ClientsDir)
	if err != nil {
		logger.Error("Directory isn't exist")
	} else {
		for _, v := range d {
			logger.Info(v.Name())
		}
	}
	switch filename {
	case "linux-GophKeeper":
		return c.File(hs.config.ClientsDir + "/linux-client")
	case "windows-GophKeeper":
		return c.File(hs.config.ClientsDir + "/windows-client")
	case "arm-linux-GophKeeper":
		return c.File(hs.config.ClientsDir + "/arm-linux-client")
	case "server.crt":
		return c.File(hs.config.CrtFile)
	default:
		return c.NoContent(http.StatusNotFound)
	}
}

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
	uid, err := GetUID(c)
	if err != nil {
		logger.Error(fmt.Errorf("Problem with uid: %w", err))
		return c.NoContent(http.StatusBadRequest)
	}

	d, err := hs.dataService.GetByUser(c.Request().Context(), uid)
	if err != nil {
		logger.Error(fmt.Errorf("Problem with GetData: %w", err))
		switch {
		case errors.Is(err, services.ErrDataNotFound):
			return c.NoContent(http.StatusNotFound)
		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, d)

}
func (hs *Server) PostData(c echo.Context) error {
	d := []datas.Data{}
	err := c.Bind(&d)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	uid, err := GetUID(c)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	logger.Debug(fmt.Sprintf("POST %d datas: %v", len(d), d))
	err = hs.dataService.Add(c.Request().Context(), uid, d)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, d)

}
func (hs *Server) DeleteData(c echo.Context) error {
	ids := []int{}
	err := c.Bind(&ids)
	if err != nil {
		logger.Error(fmt.Errorf("Invalid ids: %w", err))
		return c.NoContent(http.StatusBadRequest)
	}
	uid, err := GetUID(c)
	if err != nil {
		logger.Error(fmt.Errorf("Problem with uid: %w", err))
		return c.NoContent(http.StatusBadRequest)
	}

	err = hs.dataService.DeleteByID(c.Request().Context(), uid, ids)
	if err != nil {
		logger.Error(fmt.Errorf("Problem with GetData: %w", err))
		switch {
		case errors.Is(err, services.ErrDataNotFound):
			return c.NoContent(http.StatusNotFound)
		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.NoContent(http.StatusOK)

}
func (hs *Server) UpdateData(c echo.Context) error {
	d := []datas.Data{}
	err := c.Bind(&d)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	uid, err := GetUID(c)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	logger.Debug(fmt.Sprintf("Patch %d datas: %v", len(d), d))
	err = hs.dataService.Update(c.Request().Context(), uid, d)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, d)
}

func (hs *Server) Ping(c echo.Context) error {
	err := hs.dataService.Health(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return nil
}
