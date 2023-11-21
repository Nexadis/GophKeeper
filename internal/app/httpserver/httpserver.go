package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Nexadis/GophKeeper/internal/app/services"
	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
)

type Server struct {
	config      *config.HTTPConfig
	e           *echo.Echo
	dataService *services.Data
	authService *services.Auth
}

func New(c *config.HTTPConfig, d *services.Data, a *services.Auth) *Server {
	e := echo.New()
	hs := &Server{
		c,
		e,
		d,
		a,
	}
	hs.mountHandlers()
	return hs
}

func (hs *Server) mountHandlers() {
	hs.e.Use(middleware.Logger(),
		middleware.Gzip(),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:     true,
			LogStatus:  true,
			LogMethod:  true,
			LogLatency: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Info(
					fmt.Sprintf(
						"request: %s %d URI=%s latency=%v",
						v.Method,
						v.Status,
						v.URI,
						v.Latency,
					),
				)
				return nil
			},
		},
		),
	)
	hs.e.POST("/register", hs.Register)
	hs.e.POST("/login", hs.Login)

	apiv1 := hs.e.Group("/api/v1")

	data := apiv1.Group("/data")

	dataByID := data.Group("/:type/:id")
	dataByID.GET("", hs.GetData)
	dataByID.POST("", hs.PostData)
	dataByID.DELETE("", hs.DeleteData)
	dataByID.PATCH("", hs.UpdateData)

	dataByUser := data.Group("/user/:id")
	dataByUser.GET("", hs.GetUserData)
	dataByUser.POST("", hs.PostUserData)
	dataByUser.DELETE("", hs.DeleteUserData)
	dataByUser.PATCH("", hs.UpdateUserData)

}

func (hs *Server) Run(ctx context.Context) error {
	if hs.config.TLS {
		return hs.e.StartTLS(hs.config.Address, hs.config.CrtFile, hs.config.KeyFile)
	}
	return hs.e.Start(hs.config.Address)

}
