package service

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"zone.com/common"
)

func (s *service) statAllRes(c echo.Context) error {
	data, err := s.lgc.StatAllRes()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) statSlingByTon(c echo.Context) error {
	data, err := s.lgc.StatSlingByTon()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) getSlingUsedTop(c echo.Context) error {
	topNum, err := strconv.Atoi(c.QueryParam("topNum"))
	if err != nil {
		topNum = 10
	}
	data, err := s.lgc.GetSlingUsedTop(topNum)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) statSlingByStatus(c echo.Context) error {
	data, err := s.lgc.StatSlingByStatus()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) registerStatRoute() {
	r := s.echo.Group("/home")
	r.Use(middleware.JWTWithConfig(*s.jwtConfig))
	r.GET("/stat_all_res", s.statAllRes)
	r.GET("/stat_sling_by_ton", s.statSlingByTon)
	r.GET("/sling_used_top", s.getSlingUsedTop)
	r.GET("/stat_sling_by_status", s.statSlingByStatus)
}
