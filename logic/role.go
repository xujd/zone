package logic

import (
	"math"

	"zone.com/common"
)

// Role 用户
type Role struct {
	BaseModel
	Name   string `json:"name" gorm:"size:64"`
	Status int16  `json:"status"` // 0-正常，1-锁定，2-删除
	Remark string `json:"remark"`
}

// TableName role表
func (Role) TableName() string {
	return "t_auth_role"
}

// RoleFunc 角色权限
type RoleFunc struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	RoleID uint   `json:"roleId"`
	Funcs  string `json:"funcs"`
}

// TableName 角色权限关系表
func (RoleFunc) TableName() string {
	return "r_auth_role_func"
}

// AddRole 添加角色
func (lgc *Logics) AddRole(role *Role) error {
	// 角色名称不能为空
	if role.Name == "" {
		return common.ErrRoleNameIsNull
	}

	role0, _ := lgc.QueryRoleByName(role.Name)
	if role0 != nil {
		return common.ErrRoleAlreadyExists
	}

	if err := lgc.db.Create(&role).Error; err != nil {
		return err
	}
	return nil
}

// QueryRoleByName 查询角色
func (lgc *Logics) QueryRoleByName(name string) (*Role, error) {
	var role Role
	if err := lgc.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// QueryRoleByID 查询角色
func (lgc *Logics) QueryRoleByID(id uint) (*Role, error) {
	var role Role
	if err := lgc.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdateRole 修改角色
func (lgc *Logics) UpdateRole(role *Role) error {
	// 默认角色不准修改
	if role.ID == 1 {
		return common.ErrNoUpdate
	}

	role0, _ := lgc.QueryRoleByName(role.Name)
	if role0 != nil && role0.ID != role.ID {
		return common.ErrRoleAlreadyExists
	}

	if err := lgc.db.Save(&role).Error; err != nil {
		return err
	}
	return nil
}

// DeleteRole 删除角色
func (lgc *Logics) DeleteRole(id uint) error {
	// 根角色不准删除
	if id == 1 {
		return common.ErrNoDelete
	}
	if err := lgc.db.Where("id = ?", id).Delete(&Role{}).Error; err != nil {
		return err
	}
	return nil
}

// ListRoles 获取角色列表
func (lgc *Logics) ListRoles(name string, pageIndex int, pageSize int) (*SearchResult, error) {
	roledb := lgc.db.Model(&Role{}).Where("deleted_at IS NULL")
	if name != "" {
		roledb = lgc.db.Model(&Role{}).Where("name LIKE ?", "%"+name+"%")
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	roledb.Count(&rowCount)                                            //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var roles []Role
	if err := roledb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &roles}, nil
}

// SetUserRole 设置用户角色
func (lgc *Logics) SetUserRole(userID uint, roleIDs []uint) error {
	// 事务
	tx := lgc.db.Begin()
	// 先删除旧数据
	if err := tx.Where("user_id = ?", userID).Delete(&UserRoleRelation{}).Error; err != nil {
		return err
	}
	// 增加新关系
	for _, value := range roleIDs {
		if err := tx.Create(&UserRoleRelation{UserID: userID, RoleID: value}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// GetUserRole 查询用户角色
func (lgc *Logics) GetUserRole(userID uint) (*[]UserRoleRelation, error) {
	var userRoles []UserRoleRelation
	if err := lgc.db.Model(&UserRoleRelation{}).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	return &userRoles, nil
}

// SetRoleFuncs 设置角色权限
func (lgc *Logics) SetRoleFuncs(roleFunc *RoleFunc) error {
	// 事务
	tx := lgc.db.Begin()
	// 先删除旧数据
	if err := tx.Where("role_id = ?", roleFunc.ID).Delete(&RoleFunc{}).Error; err != nil {
		return err
	}
	// 增加新关系
	if err := tx.Create(&roleFunc).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetRoleFuncs 获取角色权限
func (lgc *Logics) GetRoleFuncs(roleID uint) (*RoleFunc, error) {
	var roleFunc RoleFunc
	if err := lgc.db.Where("role_id = ?", roleID).First(&roleFunc).Error; err != nil {
		return nil, err
	}

	return &roleFunc, nil
}
