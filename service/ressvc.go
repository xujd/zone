package service

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"zone.com/common"
	"zone.com/logic"
)

func (s *service) addSling(c echo.Context) error {
	r := new(logic.Sling)
	if err := c.Bind(r); err != nil {
		return err
	}
	// add
	if err := s.lgc.AddSling(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) updateSling(c echo.Context) error {
	r := new(logic.Sling)
	if err := c.Bind(r); err != nil {
		return err
	}
	// update
	if err := s.lgc.UpdateSling(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) deleteSling(c echo.Context) error {
	id := uint(0)
	// Sling id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// delete
	if err := s.lgc.DeleteSling(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) listSlings(c echo.Context) error {
	name := c.QueryParam("name")
	slingType, _ := strconv.Atoi(c.QueryParam("slingType"))
	maxTonnage, _ := strconv.Atoi(c.QueryParam("maxTonnage"))
	useStatus, _ := strconv.Atoi(c.QueryParam("useStatus"))
	inspectStatus, _ := strconv.Atoi(c.QueryParam("inspectStatus"))
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	// query all
	data, err := s.lgc.ListSlings(name, uint(slingType), uint(maxTonnage),
		uint(useStatus), uint(inspectStatus), pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) addCabinet(c echo.Context) error {
	r := new(logic.Cabinet)
	if err := c.Bind(r); err != nil {
		return err
	}
	// add
	if err := s.lgc.AddCabinet(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) updateCabinet(c echo.Context) error {
	r := new(logic.Cabinet)
	if err := c.Bind(r); err != nil {
		return err
	}
	// update
	if err := s.lgc.UpdateCabinet(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) deleteCabinet(c echo.Context) error {
	id := uint(0)
	// Cabinet id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// delete
	if err := s.lgc.DeleteCabinet(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

func (s *service) listCabinets(c echo.Context) error {
	name := c.QueryParam("name")
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	// query all
	data, err := s.lgc.ListCabinets(name, pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) listGrids(c echo.Context) error {
	id := uint(0)
	// Cabinet id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// query
	data, err := s.lgc.ListCabinetGrids(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) registerResRoute() {
	r := s.echo.Group("/res")
	r.Use(middleware.JWTWithConfig(*s.jwtConfig))
	// sling
	r.POST("/sling", s.addSling)
	r.PUT("/sling", s.updateSling)
	r.DELETE("/sling/:id", s.deleteSling)
	r.GET("/slings", s.listSlings)
	// cabinet
	r.POST("/cabinet", s.addCabinet)
	r.PUT("/cabinet", s.updateCabinet)
	r.DELETE("/cabinet/:id", s.deleteCabinet)
	r.GET("/cabinets", s.listCabinets)
	r.GET("/cabinet_grids/:id", s.listGrids)
}
