package logic

import (
	"crypto/sha256"
	"fmt"
	"math"
	"time"

	"zone.com/common"
)

// User 用户
type User struct {
	BaseModel
	Name      string    `json:"name" gorm:"size:64"`
	Password  string    `json:"password" gorm:"size:64"`
	StartTime *JSONTime `json:"startTime" gorm:"type:timestamp"`
	EndTime   *JSONTime `json:"endTime" gorm:"type:timestamp"`
	Status    int16     `json:"status"` // 0-正常，1-锁定，2-删除
	Remark    string    `json:"remark"`
	StaffName string    `json:"staffName" gorm:"-"`
	StaffID   uint      `json:"staffId"`
}

// TableName user表
func (User) TableName() string {
	return "t_auth_user"
}

func (user *User) Check() error {
	// 用户状态异常
	if user.Status != 0 {
		return common.ErrUserStatus
	}

	timeFormatStr := "2006-01-02 15:04:05"
	// 开始生效时间
	if user.StartTime != nil {
		startTime := (*time.Time)(user.StartTime).Format(timeFormatStr)
		t1, _ := time.ParseInLocation(timeFormatStr, string(startTime), time.Local)
		if t1.After(time.Now()) {
			return common.ErrUserNotEffective
		}
	}
	// 结束生效时间
	if user.EndTime != nil {
		endTime := (*time.Time)(user.EndTime).Format(timeFormatStr)
		t2, _ := time.ParseInLocation(timeFormatStr, string(endTime), time.Local)
		if t2.Before(time.Now()) {
			return common.ErrUserExpired
		}
	}

	return nil
}

// UserInfo 用户信息
type UserInfo struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
	ID           uint     `json:"id"`
	StaffName    string   `json:"staffName"`
}

// UserRoleRelation 用户角色关系
type UserRoleRelation struct {
	ID     uint `json:"id" gorm:"primary_key"`
	UserID uint `json:"userId"`
	RoleID uint `json:"roleId"`
}

// TableName user_role表
func (UserRoleRelation) TableName() string {
	return "r_auth_user_role"
}

// QueryUserByName 查询用户
func (lgc *Logics) QueryUserByName(name string) (*User, error) {

	var user User
	if err := lgc.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, common.ErrUserNotFound
	}

	return &user, nil
}

// QueryUserByID 查询用户
func (lgc *Logics) QueryUserByID(id uint) (*User, error) {

	var user User
	selectStr := "t_auth_user.id,t_auth_user.created_at,t_auth_user.updated_at,t_auth_user.deleted_at,t_auth_user.name,t_auth_user.start_time,t_auth_user.end_time,t_auth_user.status,t_auth_user.remark,t_auth_user.staff_id, t_sys_staff.name AS staff_name"

	if err := lgc.db.Table("t_auth_user").Select(selectStr).
		Joins("JOIN t_sys_staff ON t_auth_user.staff_id = t_sys_staff.id").
		Where("t_auth_user.deleted_at IS NULL AND t_auth_user.id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// AddUser 添加用户
func (lgc *Logics) AddUser(user *User) error {
	// 默认密码
	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte("123456a?"+user.Name)))
	// 用户名不能为空
	if user.Name == "" {
		return common.ErrUserNameIsNull
	}
	// 员工未指定
	if user.StaffID == 0 {
		return common.ErrUserStaffIsNull
	}

	user0, _ := lgc.QueryUserByName(user.Name)
	if user0 != nil {
		return common.ErrUserAlreadyExists
	}

	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password+user.Name)))
	if err := lgc.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUser 修改用户
func (lgc *Logics) UpdateUser(user *User) error {
	// 默认用户不准修改
	if user.ID == 1 {
		return common.ErrNoUpdate
	}

	// 用户名不能为空
	if user.Name == "" {
		return common.ErrUserNameIsNull
	}
	// 员工未指定
	if user.StaffID == 0 {
		return common.ErrUserStaffIsNull
	}
	user0, _ := lgc.QueryUserByName(user.Name)
	if user0 != nil && user0.ID != user.ID {
		return common.ErrUserAlreadyExists
	}
	data := map[string]interface{}{
		"Remark":  user.Remark,
		"StaffID": user.StaffID,
	}
	if user.StartTime != nil {
		data["StartTime"] = user.StartTime
	}
	if user.EndTime != nil {
		data["EndTime"] = user.EndTime
	}
	if user.Status > -1 {
		data["Status"] = user.Status
	}
	if err := lgc.db.Model(&user).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser 删除用户
func (lgc *Logics) DeleteUser(id uint) error {
	// 根用户不准删除
	if id == 1 {
		return common.ErrNoDelete
	}
	if err := lgc.db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// ListUsers 查询用户
func (lgc *Logics) ListUsers(name string, pageIndex int, pageSize int) (*SearchResult, error) {

	selectStr := "t_auth_user.id,t_auth_user.created_at,t_auth_user.updated_at,t_auth_user.deleted_at,t_auth_user.name,t_auth_user.start_time,t_auth_user.end_time,t_auth_user.status,t_auth_user.remark,t_auth_user.staff_id, t_sys_staff.name AS staff_name"
	userdb := lgc.db.Table("t_auth_user").Select(selectStr).
		Joins("JOIN t_sys_staff ON t_auth_user.staff_id = t_sys_staff.id").
		Where("t_auth_user.deleted_at IS NULL")

	if name != "" {
		userdb = userdb.Where("t_auth_user.name LIKE ?", "%"+name+"%")
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	var rowCount int64
	userdb.Count(&rowCount)                                            //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var users []User
	if err := userdb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &users}, nil
}

// GetUserInfo 用户信息
func (lgc *Logics) GetUserInfo(userId uint) (*UserInfo, error) {
	user, err := lgc.QueryUserByID(userId)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		Introduction: user.Remark,
		Avatar:       "./assets/user.gif",
		Name:         user.Name,
		ID:           user.ID,
		StaffName:    user.StaffName,
	}

	var roles []Role
	lgc.db.Raw("SELECT * FROM t_auth_role WHERE id IN (SELECT role_id FROM r_auth_user_role WHERE user_id=?)", user.ID).Scan(&roles)

	userInfo.Roles = make([]string, len(roles))
	for key, value := range roles {
		userInfo.Roles[key] = value.Name
	}

	return userInfo, nil
}

// ResetPassword 重置密码
func (lgc *Logics) ResetPassword(userID uint) (string, error) {
	// 默认用户不准修改
	if userID == 1 {
		return "", common.ErrNoUpdate
	}

	user0, err0 := lgc.QueryUserByID(userID)
	if err0 != nil {
		return "", common.ErrUserNotFound
	}
	// 两次加密
	password1 := fmt.Sprintf("%x", sha256.Sum256([]byte("123456a?"+user0.Name)))
	data := map[string]interface{}{
		"Password": fmt.Sprintf("%x", sha256.Sum256([]byte(password1+user0.Name))),
	}
	if err := lgc.db.Model(&user0).Updates(data).Error; err != nil {
		return "", err
	}
	return "success", nil
}

// UpdatePassword 修改密码
func (lgc *Logics) UpdatePassword(userID uint, password string, newPassword string) (string, error) {
	// 默认用户不准修改
	if userID == 1 {
		return "", common.ErrNoUpdate
	}

	// 查询用户是否存在
	var user User
	if err := lgc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	// 不存在
	if user.ID == 0 {
		return "", common.ErrUserNotFound
	}
	// 确认密码
	oldPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password+user.Name)))
	if oldPassword != user.Password {
		return "", common.ErrPwdDismatch
	}
	data := map[string]interface{}{
		"Password": fmt.Sprintf("%x", sha256.Sum256([]byte(newPassword+user.Name))),
	}
	if err := lgc.db.Model(&user).Updates(data).Error; err != nil {
		return "", err
	}
	return "success", nil
}
