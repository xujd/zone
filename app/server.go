package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"zone.com/common"
	"zone.com/service"
	"zone.com/util"
)

func Run(op *ServerOption) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Enable metrics middleware
	p := prometheus.NewPrometheus("zone", nil)
	p.Use(e)
	// rate limit
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	// welcome
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Welcome to ZPPR!")
	})
	// error handler
	e.HTTPErrorHandler = httpErrorHandler

	// static files directory
	util.FileDir = op.FileDir
	if !util.HasSuffix(util.FileDir, "/") {
		util.FileDir = util.FileDir + "/"
	}

	// db
	db, err := util.InitDB(op.DbHost, op.DbUser, op.DbPassword, op.DbName, op.DbPort)
	if err != nil {
		e.Logger.Fatal("DB init failed.")
		return err
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	// Service
	svc := service.NewService()
	svc.SetConfig(e, db)
	svc.RegisterServices()

	// Start server
	go func() {
		if err := e.Start(op.AddrPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
		return err
	}
	return nil
}

func httpErrorHandler(err error, c echo.Context) {
	var (
		code    = common.ERR_BAD_REQUEST
		success = false
		message = fmt.Sprint(err)
	)

	if e, ok := err.(*echo.HTTPError); ok {
		code = codeFormat(e.Code)
		message = fmt.Sprint(e.Message)
	}
	if he, ok := err.(*common.HttpError); ok {
		code = he.Code
		message = he.Message
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			err := c.JSON(http.StatusOK, common.NewHttpMsg(code, success, message, nil))
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

func codeFormat(code int) int {
	switch code {
	case http.StatusUnauthorized:
		return common.ERR_ILLEGAL_TOKEN
	default:
		return common.ERR_INTERNAL_SERVER_ERROR
	}
}
