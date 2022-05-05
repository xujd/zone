package logic

import (
	"fmt"
	"io"
	"math"
	"os"

	"zone.com/common"
	"zone.com/util"
)

// Company 公司
type Company struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Name   string `json:"name" gorm:"size:128"` // 公司名称
	Status int16  `json:"status"`               // 状态：0-正常，1-停用
	Remark string `json:"remark"`               // 说明
}

// TableName company表
func (Company) TableName() string {
	return "t_sys_company"
}

// Department 部门
type Department struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Name      string  `json:"name" gorm:"size:128"`                // 部门名称
	Company   Company `json:"company" gorm:"ForeignKey:CompanyID"` // 公司
	CompanyID uint    `json:"companyId"`                           // 公司ID
	Status    int16   `json:"status"`                              // 状态：0-正常，1-停用
	Remark    string  `json:"remark"`                              // 说明
}

// TableName department
func (Department) TableName() string {
	return "t_sys_department"
}

// Staff 员工
type Staff struct {
	BaseModel
	Name           string    `json:"name" gorm:"size:64"`            // 员工姓名
	CompanyName    string    `json:"companyName" gorm:"-"`           // 公司
	CompanyID      uint      `json:"companyId"`                      // 公司ID
	DepartmentName string    `json:"departmentName" gorm:"-"`        // 部门
	DepartmentID   uint      `json:"departmentId"`                   // 部门ID
	PostName       string    `json:"postName"`                       // 职务
	Birthday       *JSONTime `json:"birthday" gorm:"type:timestamp"` // 出生日期
	Status         int16     `json:"status"`                         // 状态：0-正常
	Remark         string    `json:"remark"`                         // 说明
}

// TableName staff表
func (Staff) TableName() string {
	return "t_sys_staff"
}

// DictData 字典数据
type DictData struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Key   int    `json:"key"`
	Name  string `json:"name" gorm:"size:32"`  // 名称
	Type  string `json:"type" gorm:"size:32"`  // 类型
	Note  string `json:"note" gorm:"size:64"`  // 备注
	Scene string `json:"scene" gorm:"size:32"` // 应用场景
}

// TableName DictData
func (DictData) TableName() string {
	return "t_sys_dict"
}

// ListCompanys 查询公司
func (lgc *Logics) ListCompanys(name string, pageIndex int, pageSize int) (*SearchResult, error) {
	companydb := lgc.db.Model(&Company{})
	if name != "" {
		companydb = lgc.db.Model(&Company{}).Where("name LIKE ?", "%"+name+"%")
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	companydb.Count(&rowCount)                                         //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var companys []Company
	if err := companydb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&companys).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &companys}, nil
}

// ListDepartments 查询部门
func (lgc *Logics) ListDepartments(name string, companyID uint, pageIndex int, pageSize int) (*SearchResult, error) {
	deptdb := lgc.db.Model(&Department{})
	if name != "" {
		deptdb = lgc.db.Model(&Department{}).Where("name LIKE ?", "%"+name+"%")
	}
	if companyID > 0 {
		deptdb = deptdb.Where("company_id = ?", companyID)
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	deptdb.Count(&rowCount)                                            //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var deptList []Department
	if err := deptdb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&deptList).Error; err != nil {
		return nil, err
	}

	// 关联公司
	for key, dept := range deptList {
		lgc.db.Model(&dept).Association("Company").Find(&dept.Company)
		deptList[key] = dept
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &deptList}, nil
}

// AddStaff 添加员工
func (lgc *Logics) AddStaff(staff *Staff) error {
	// 员工姓名不能为空
	if staff.Name == "" {
		return common.ErrStaffNameIsNull
	}

	if err := lgc.db.Create(&staff).Error; err != nil {
		return err
	}

	// 员工的照片
	srcFile := "./temp/temp.jpg"
	srcDefault := "./temp/default.jpg"
	fileName := fmt.Sprintf("./temp/%06d.jpg", staff.ID)
	if fileExist(srcFile) {
		os.Rename(srcFile, fileName)
		copy(fileName, fmt.Sprintf("%s%06d.jpg", util.FileDir, staff.ID))
		// 覆盖掉web目录的temp.jpg
		copy(srcDefault, fmt.Sprintf("%stemp.jpg", util.FileDir))
	}
	return nil
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// UpdateStaff 修改员工
func (lgc *Logics) UpdateStaff(staff *Staff) error {
	// 默认员工不准修改
	if staff.ID == 1 {
		return common.ErrNoUpdate
	}

	// 员工姓名不能为空
	if staff.Name == "" {
		return common.ErrStaffNameIsNull
	}
	if err := lgc.db.Save(&staff).Error; err != nil {
		return err
	}
	return nil
}

// DeleteStaff 删除员工
func (lgc *Logics) DeleteStaff(id uint) error {
	// 默认员工不准删除
	if id == 1 {
		return common.ErrNoDelete
	}
	if err := lgc.db.Where("id = ?", id).Delete(&Staff{}).Error; err != nil {
		return err
	}
	return nil
}

// QueryStaffByID 查询员工
func (lgc *Logics) QueryStaffByID(id uint) (*Staff, error) {

	var staff Staff
	if err := lgc.db.Where("id = ?", id).First(&staff).Error; err != nil {
		return nil, err
	}

	return &staff, nil
}

// ListStaffs 查询员工
func (lgc *Logics) ListStaffs(name string, companyID uint, departmentID uint, pageIndex int, pageSize int) (*SearchResult, error) {

	staffdb := lgc.db.Table("t_sys_staff").
		Select("t_sys_staff.*, t_sys_company.name AS company_name, t_sys_department.name AS department_name").
		Joins("JOIN t_sys_company ON t_sys_staff.company_id = t_sys_company.id").
		Joins("JOIN t_sys_department ON t_sys_staff.department_id = t_sys_department.id").
		Where("t_sys_staff.deleted_at IS NULL")

	if name != "" {
		staffdb = staffdb.Where("t_sys_staff.name LIKE ?", "%"+name+"%")
	}

	if companyID > 0 {
		staffdb = staffdb.Where("t_sys_staff.company_id = ?", companyID)
	}

	if departmentID > 0 {
		staffdb = staffdb.Where("t_sys_staff.department_id = ?", departmentID)
	}

	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	staffdb.Count(&rowCount)                                           //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var staffs []Staff
	if err := staffdb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&staffs).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &staffs}, nil
}

// ListDict 查询字典
func (lgc *Logics) ListDict(scene string, dictType string) (*[]DictData, error) {

	var dictDatas []DictData

	db := lgc.db.Model(&DictData{})

	if scene != "" {
		db = db.Where("scene = ?", scene)
	}
	if dictType != "" {
		db = db.Where("type = ?", dictType)
	}
	if err := db.Find(&dictDatas).Error; err != nil {
		return nil, err
	}

	return &dictDatas, nil
}

// 添加字典
func (lgc *Logics) AddDict(dict *DictData) error {
	if err := lgc.db.Create(&dict).Error; err != nil {
		return err
	}
	return nil
}

// 修改字典
func (lgc *Logics) UpdateDict(dict *DictData) error {
	if err := lgc.db.Save(&dict).Error; err != nil {
		return err
	}
	return nil
}

// 删除字典
func (lgc *Logics) DeleteDict(id uint) error {
	if err := lgc.db.Where("id = ?", id).Delete(&DictData{}).Error; err != nil {
		return err
	}
	return nil
}
