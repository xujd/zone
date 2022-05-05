package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"zone.com/common"
	"zone.com/logic"
)

func (s *service) store(c echo.Context) error {
	req, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	if err := json.Unmarshal(req, &data); err != nil {
		return err
	}
	if data["cabinetId"] == nil || data["cabinetId"] == "" ||
		data["gridNo"] == nil || data["gridNo"] == "" ||
		data["resId"] == nil || data["resId"] == "" {
		return common.ErrBadQueryParams
	}
	cabinetId, _ := strconv.Atoi(data["cabinetId"].(string))
	gridNo, _ := strconv.Atoi(data["gridNo"].(string))
	resId, _ := strconv.Atoi(data["resId"].(string))
	if err := s.lgc.Store(uint(cabinetId), uint(gridNo), uint(resId)); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) takeReturn(c echo.Context) error {
	req, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	if err := json.Unmarshal(req, &data); err != nil {
		return err
	}
	if data["cabinetId"] == nil || data["cabinetId"] == "" ||
		data["gridNo"] == nil || data["gridNo"] == "" ||
		data["flag"] == nil || data["flag"] == "" {
		return common.ErrBadQueryParams
	}
	cabinetId, _ := strconv.Atoi(data["cabinetId"].(string))
	gridNo, _ := strconv.Atoi(data["gridNo"].(string))
	flag, _ := strconv.Atoi(data["flag"].(string))
	if err := s.lgc.TakeReturn(uint(cabinetId), uint(gridNo), flag); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) takeReturnByResID(c echo.Context) error {
	u := new(logic.UseLog)
	if err := c.Bind(u); err != nil {
		return err
	}
	// do
	if err := s.lgc.TakeReturnByResID(u); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) getResUseLog(c echo.Context) error {
	param := &logic.UseLogQueryParam{
		ResName:       c.QueryParam("resName"),
		TakeStartTime: c.QueryParam("takeStartTime"),
		TakeEndTime:   c.QueryParam("takeEndTime"),
	}
	returnFlag, _ := strconv.Atoi(c.QueryParam("returnFlag"))
	takeStaff, _ := strconv.Atoi(c.QueryParam("takeStaff"))
	returnStaff, _ := strconv.Atoi(c.QueryParam("returnStaff"))
	param.ReturnFlag = returnFlag
	param.TakeStaff = uint(takeStaff)
	param.ReturnStaff = uint(returnStaff)

	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	data, err := s.lgc.GetTakeReturnLog(param, pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) registerUsageRoute() {
	r := s.echo.Group("/res")
	r.Use(middleware.JWTWithConfig(*s.jwtConfig))
	// usage
	r.POST("/store", s.store)
	r.POST("/take_return", s.takeReturn)
	r.POST("/take_return_by_res", s.takeReturnByResID)
	r.GET("/uselog", s.getResUseLog)
}
