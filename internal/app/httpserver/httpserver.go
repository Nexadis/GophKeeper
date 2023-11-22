package httpserver

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

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
	apiv1.GET("/ping", hs.Ping)

	data := apiv1.Group("/data")

	data.GET("", hs.GetData)
	data.POST("", hs.PostData)
	data.DELETE("", hs.DeleteData)
	data.PATCH("", hs.UpdateData)

	dataByUser := apiv1.Group("/user")
	dataByUser.GET("", hs.GetUserData)
	dataByUser.POST("", hs.PostUserData)
	dataByUser.DELETE("", hs.DeleteUserData)
	dataByUser.PATCH("", hs.UpdateUserData)

}

func (hs *Server) Run(ctx context.Context) error {
	hs.mustSecret()
	if hs.config.TLS {
		return hs.e.StartTLS(hs.config.Address, hs.config.CrtFile, hs.config.KeyFile)
	}
	logger.Info("Server is running without TLS. Be careful, your password may be intercepted!")
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		err := hs.e.Start(hs.config.Address)
		if err != nil {
			logger.Error()
		}

	}()

	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return hs.e.Shutdown(ctx)

}

func (hs *Server) mustSecret() {
	if hs.config.JWTSecret == nil {
		buf := make([]byte, 32)
		_, err := rand.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		hs.config.JWTSecret = buf
	}

}
