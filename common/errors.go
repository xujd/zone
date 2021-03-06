package common

import "errors"

var (
	// ErrAlreadyExists 记录已存在
	ErrAlreadyExists = errors.New("记录已存在")
	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("用户已存在")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")
	// ErrNameOrPasswordIsNull 用户名或密码为空
	ErrNameOrPasswordIsNull = errors.New("用户名或密码为空")
	// ErrUserNotEffective 用户还未生效
	ErrUserNotEffective = errors.New("用户还未生效")
	// ErrUserExpired 用户已失效
	ErrUserExpired = errors.New("用户已失效")
	// ErrUserStaffIsNull 用户员工未选择
	ErrUserStaffIsNull = errors.New("用户员工未选择")
	// ErrUserNameIsNull 用户名为空
	ErrUserNameIsNull = errors.New("用户名为空")
	// ErrNotFound 未找到数据表或相关记录
	ErrNotFound = errors.New("未找到数据表或相关记录")
	// ErrUserPwdDismatch 用户名或密码错误
	ErrUserPwdDismatch = errors.New("用户名或密码错误")
	// ErrPwdDismatch 密码错误
	ErrPwdDismatch = errors.New("密码错误")
	// ErrUserStatus 用户状态异常
	ErrUserStatus = errors.New("用户状态异常")
	// ErrBadQueryParams 查询参数错误
	ErrBadQueryParams = errors.New("查询参数错误")
	// ErrRoleNameIsNull 角色名称为空
	ErrRoleNameIsNull = errors.New("角色名称为空")
	// ErrRoleAlreadyExists 角色已存在
	ErrRoleAlreadyExists = errors.New("角色已存在")
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("角色不存在")
	// ErrStaffNameIsNull 员工姓名不能为空
	ErrStaffNameIsNull = errors.New("员工姓名不能为空")
	// ErrStaddIDAlreadyExists 员工编号已存在
	ErrStaddIDAlreadyExists = errors.New("员工编号已存在")
	// ErrNoDelete 数据不允许删除
	ErrNoDelete = errors.New("数据不允许删除")
	// ErrNoUpdate 数据不允许修改
	ErrNoUpdate = errors.New("数据不允许修改")
	// ErrSlingNameIsNull 吊索具名称不能为空
	ErrSlingNameIsNull = errors.New("吊索具名称不能为空")
	// ErrSlingRfIDIsNull 吊索具RFID不能为空
	ErrSlingRfIDIsNull = errors.New("吊索具RFID不能为空")
	// ErrSlingAlreadyExists 吊索具名称或RFID重复
	ErrSlingAlreadyExists = errors.New("吊索具名称或RFID重复")
	// ErrSlingCabinetIsNull 吊索具存放位置不能为空
	ErrSlingCabinetIsNull = errors.New("吊索具存放位置不能为空")
	// ErrCabinetNameIsNull 智能柜名称不能为空
	ErrCabinetNameIsNull = errors.New("智能柜名称不能为空")
	// ErrCabinetAlreadyExists 智能柜名称重复
	ErrCabinetAlreadyExists = errors.New("智能柜名称重复")
	// ErrCabinetGridIsZero 智能柜箱格数不能是0
	ErrCabinetGridIsZero = errors.New("智能柜箱格数不能是0")
	// ErrGridAlreadyInUse 箱格已使用
	ErrGridAlreadyInUse = errors.New("箱格已使用")
)
