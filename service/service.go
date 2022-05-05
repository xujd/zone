package service

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"zone.com/logic"
	"zone.com/util"
)

// Service service methods
type Service interface {
	SetConfig(echo *echo.Echo, db *gorm.DB)
	RegisterServices()
}

// NewService create a new service instance
func NewService() Service {
	return new(service)
}

type service struct {
	echo      *echo.Echo
	db        *gorm.DB
	lgc       *logic.Logics
	jwtConfig *middleware.JWTConfig
}

// SetConfig
func (s *service) SetConfig(echo *echo.Echo, db *gorm.DB) {
	s.echo = echo
	s.db = db
	s.lgc = logic.NewLogics(db)
	// Configure middleware with the custom claims type
	s.jwtConfig = &middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(util.SecretKey),
	}
}

// RegisterServices
func (s *service) RegisterServices() {
	// auth
	s.registerAuthRoute()
	// home
	s.registerStatRoute()
	// sys
	s.registerSysRoute()
	// res
	s.registerResRoute()
	// usage
	s.registerUsageRoute()

	// file upload
	r := s.echo.Group("/file")
	r.Use(middleware.JWTWithConfig(*s.jwtConfig))
	r.POST("/upload", s.upload)

}
