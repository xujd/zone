package service

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"zone.com/common"
	"zone.com/logic"
)

func (s *service) listCompanys(c echo.Context) error {
	name := c.QueryParam("name")
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	data, err := s.lgc.ListCompanys(name, pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) listDepartments(c echo.Context) error {
	name := c.QueryParam("name")
	companyId, _ := strconv.Atoi(c.QueryParam("companyId"))
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	data, err := s.lgc.ListDepartments(name, uint(companyId), pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) addStaff(c echo.Context) error {
	r := new(logic.Staff)
	if err := c.Bind(r); err != nil {
		return err
	}
	// add
	if err := s.lgc.AddStaff(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) updateStaff(c echo.Context) error {
	r := new(logic.Staff)
	if err := c.Bind(r); err != nil {
		return err
	}
	// update
	if err := s.lgc.UpdateStaff(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) deleteStaff(c echo.Context) error {
	id := uint(0)
	// staff id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// delete
	if err := s.lgc.DeleteStaff(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) listStaffs(c echo.Context) error {
	name := c.QueryParam("name")
	companyId, _ := strconv.Atoi(c.QueryParam("companyId"))
	departmentId, _ := strconv.Atoi(c.QueryParam("departmentId"))
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	// query all
	data, err := s.lgc.ListStaffs(name, uint(companyId), uint(departmentId), pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) addDict(c echo.Context) error {
	r := new(logic.DictData)
	if err := c.Bind(r); err != nil {
		return err
	}
	// add
	if err := s.lgc.AddDict(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) updateDict(c echo.Context) error {
	r := new(logic.DictData)
	if err := c.Bind(r); err != nil {
		return err
	}
	// update
	if err := s.lgc.UpdateDict(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) deleteDict(c echo.Context) error {
	id := uint(0)
	// dict id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// delete
	if err := s.lgc.DeleteDict(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) listDict(c echo.Context) error {
	scene := c.QueryParam("scene")
	dictType := c.QueryParam("type")
	// query all
	data, err := s.lgc.ListDict(scene, dictType)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}
func (s *service) registerSysRoute() {
	r := s.echo.Group("/sys")
	r.Use(middleware.JWTWithConfig(*s.jwtConfig))
	r.GET("/companys", s.listCompanys)
	r.GET("/departments", s.listDepartments)
	// staff
	r.POST("/staff", s.addStaff)
	r.PUT("/staff", s.updateStaff)
	r.DELETE("/staff/:id", s.deleteStaff)
	r.GET("/staffs", s.listStaffs)
	// dict
	r.POST("/dict", s.addDict)
	r.PUT("/dict", s.updateDict)
	r.DELETE("/dict/:id", s.deleteDict)
	r.GET("/dict", s.listDict)
}
