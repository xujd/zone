package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"zone.com/common"
	"zone.com/logic"
	"zone.com/util"
)

type jwtCustomClaims struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Admin  bool   `json:"admin"`
	jwt.StandardClaims
}

func (s *service) login(c echo.Context) error {
	u := new(logic.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// form value
	if u.Name == "" {
		u.Name = c.FormValue("name")
		u.Password = c.FormValue("password")
	}

	// error
	if u.Name == "" || u.Password == "" {
		return common.ErrBadQueryParams
	}

	user, err := s.lgc.QueryUserByName(u.Name)
	if err != nil {
		return err
	}
	// 用户状态异常
	if err := user.Check(); err != nil {
		return err
	}

	passwordNew := fmt.Sprintf("%x", sha256.Sum256([]byte(u.Password+u.Name)))

	// Throws unauthorized error
	if user.Name != u.Name || user.Password != passwordNew {
		return common.ErrUserPwdDismatch
	}

	t, err := s.signToken(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData(echo.Map{
		"token": t,
	}))
}

func (s *service) signToken(user *logic.User) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		fmt.Sprint(user.ID),
		user.Name,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(util.SecretKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *service) renewval(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	username := claims.Name
	u, err := s.lgc.QueryUserByName(username)
	if err != nil {
		s.echo.Logger.Error(err)
		return err
	}
	// 用户状态异常
	if err := u.Check(); err != nil {
		return err
	}

	t, err := s.signToken(u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData(echo.Map{
		"token": t,
	}))
}

// logout 退出登录
func (s *service) logout(c echo.Context) error {
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// addUser
func (s *service) addUser(c echo.Context) error {
	u := new(logic.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// add
	if err := s.lgc.AddUser(u); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// updateUser
func (s *service) updateUser(c echo.Context) error {
	u := new(logic.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// update
	if err := s.lgc.UpdateUser(u); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// deleteUser
func (s *service) deleteUser(c echo.Context) error {
	id := uint(0)
	// user id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// delete
	if err := s.lgc.DeleteUser(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// queryUserByID
func (s *service) queryUserByID(c echo.Context) error {
	id := uint(0)
	// user id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// query by id
	user, err := s.lgc.QueryUserByID(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(user))
}

// listUsers
func (s *service) listUsers(c echo.Context) error {
	name := c.QueryParam("name")
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	// query all
	user, err := s.lgc.ListUsers(name, pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(user))
}

// getUserInfo
func (s *service) getUserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	userId, err := strconv.Atoi(claims.UserId)
	if err != nil {
		return err
	}
	userInfo, err := s.lgc.GetUserInfo(uint(userId))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(userInfo))
}

// resetPassword
func (s *service) resetPassword(c echo.Context) error {
	req, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	if err := json.Unmarshal(req, &data); err != nil {
		return err
	}
	if data["userId"] == nil || data["userId"] == "" {
		return common.ErrBadQueryParams
	}

	if _, err := s.lgc.ResetPassword(uint(data["userId"].(float64))); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// updatePassword
func (s *service) updatePassword(c echo.Context) error {
	req, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	if err := json.Unmarshal(req, &data); err != nil {
		return err
	}
	if data["userId"] == nil || data["userId"] == "" ||
		data["password"] == nil || data["password"] == "" ||
		data["newPassword"] == nil || data["newPassword"] == "" {
		return common.ErrBadQueryParams
	}

	if _, err := s.lgc.UpdatePassword(uint(data["userId"].(float64)), data["password"].(string), data["newPassword"].(string)); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// addRole
func (s *service) addRole(c echo.Context) error {
	r := new(logic.Role)
	if err := c.Bind(r); err != nil {
		return err
	}
	// add
	if err := s.lgc.AddRole(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// updateRole
func (s *service) updateRole(c echo.Context) error {
	r := new(logic.Role)
	if err := c.Bind(r); err != nil {
		return err
	}
	// update
	if err := s.lgc.UpdateRole(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// deleteRole
func (s *service) deleteRole(c echo.Context) error {
	id := uint(0)
	// role id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// delete
	if err := s.lgc.DeleteRole(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// queryRoleByID
func (s *service) queryRoleByID(c echo.Context) error {
	id := uint(0)
	// role id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// query by id
	role, err := s.lgc.QueryRoleByID(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(role))
}

// listRoles
func (s *service) listRoles(c echo.Context) error {
	name := c.QueryParam("name")
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	// query all
	roles, err := s.lgc.ListRoles(name, pageIndex, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(roles))
}

// setUserRole
func (s *service) setUserRole(c echo.Context) error {
	req, _ := ioutil.ReadAll(c.Request().Body)
	var data map[string]interface{}
	if err := json.Unmarshal(req, &data); err != nil {
		return err
	}
	roleIds := make([]uint, len(data["roleIds"].([]interface{})))
	for i, value := range data["roleIds"].([]interface{}) {
		roleIds[i] = uint(value.(float64))
	}
	if err := s.lgc.SetUserRole(uint(data["userId"].(float64)), roleIds); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// getUserRole
func (s *service) getUserRole(c echo.Context) error {
	id := uint(0)
	// id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// query by id
	data, err := s.lgc.GetUserRole(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

// setRoleFunc
func (s *service) setRoleFunc(c echo.Context) error {
	r := new(logic.RoleFunc)
	if err := c.Bind(r); err != nil {
		return err
	}
	// add
	if err := s.lgc.SetRoleFuncs(r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData("success"))
}

// getRoleFunc
func (s *service) getRoleFunc(c echo.Context) error {
	id := uint(0)
	// id
	if err := echo.PathParamsBinder(c).Uint("id", &id).BindError(); err != nil {
		return err
	}
	// query by id
	data, err := s.lgc.GetRoleFuncs(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, common.NewHttpMsgData(data))
}

func (s *service) registerAuthRoute() {
	// login
	s.echo.POST("/login", s.login)
	// logout
	s.echo.POST("/logout", s.logout)

	r := s.echo.Group("/auth")
	r.Use(middleware.JWTWithConfig(*s.jwtConfig))
	r.GET("/renewval", s.renewval)
	// user
	r.POST("/user", s.addUser)
	r.PUT("/user", s.updateUser)
	r.DELETE("/user/:id", s.deleteUser)
	r.GET("/user/:id", s.queryUserByID)
	r.GET("/users", s.listUsers)
	r.GET("/userinfo", s.getUserInfo)
	r.POST("/resetpwd", s.resetPassword)
	r.POST("/updatepwd", s.updatePassword)
	// role
	r.POST("/role", s.addRole)
	r.PUT("/role", s.updateRole)
	r.DELETE("/role/:id", s.deleteRole)
	r.GET("/role/:id", s.queryRoleByID)
	r.GET("/roles", s.listRoles)
	r.POST("/userrole", s.setUserRole)
	r.GET("/userrole/:id", s.getUserRole)
	r.POST("/rolefunc", s.setRoleFunc)
	r.GET("/rolefunc/:id", s.getRoleFunc)
}
