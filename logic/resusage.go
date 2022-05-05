package logic

import (
	"math"

	"zone.com/common"
)

// UseLog 使用日志
type UseLog struct {
	BaseModel
	ResID           uint      `json:"resId"`
	Flag            int       `json:"flag" gorm:"-"` // 借还标记 0-还，1-借
	RfID            string    `json:"rfId" gorm:"size:64"`
	ResName         string    `json:"resName" gorm:"size:64"` // 资产名称
	TakeStaffID     uint      `json:"takeStaffId"`
	TakeStaffName   string    `json:"takeStaffName" gorm:"size:64"`         // 借用人姓名
	TakeTime        *JSONTime `json:"takeTime" gorm:"type:timestamp"`       // 借用时间
	ReturnPlanTime  JSONTime  `json:"returnPlanTime" gorm:"type:timestamp"` // 预计归还时间
	ReturnStaffID   uint      `json:"returnStaffId"`
	ReturnStaffName string    `json:"returnStaffName" gorm:"size:64"`   // 归还人姓名
	ReturnTime      *JSONTime `json:"returnTime" gorm:"type:timestamp"` // 归还时间
	Remark          string    `json:"remark"`                           // 说明
}

// TableName UseLog
func (UseLog) TableName() string {
	return "t_res_use_log"
}

type UseLogQueryParam struct {
	ResName       string
	ReturnFlag    int
	TakeStaff     uint
	ReturnStaff   uint
	TakeStartTime string
	TakeEndTime   string
}

// Store 存
func (lgc *Logics) Store(cabinetID uint, gridNo uint, resID uint) error {
	// 判重
	var cabinetGrid CabinetGrid
	lgc.db.Model(&CabinetGrid{}).Where("cabinet_id = ? and grid_no = ?", cabinetID, gridNo).First(&cabinetGrid)
	if cabinetGrid.InResID > 0 && cabinetGrid.InResID != resID {
		return common.ErrGridAlreadyInUse
	}
	// 事务
	tx := lgc.db.Begin()
	// 是否存在
	var cabinetGrid0 CabinetGrid
	if err0 := tx.Where("in_res_id = ?", resID).First(&cabinetGrid0).Error; err0 == nil {
		// 更新
		if err := tx.Model(&cabinetGrid0).Updates(CabinetGrid{CabinetID: cabinetID, GridNo: gridNo}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 创建新纪录
		if err := tx.Create(&CabinetGrid{GridNo: gridNo, CabinetID: cabinetID, InResID: resID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// TakeReturn 取-将is_out设置为1;还-将is_out设置为0
func (lgc *Logics) TakeReturn(cabinetID uint, gridNo uint, flag int) error {

	if err := lgc.db.Model(&CabinetGrid{}).Where("cabinet_id = ? and grid_no = ?", cabinetID, gridNo).Update("is_out", flag).Error; err != nil {
		return err
	}
	return nil
}

// TakeReturnByResID 按资源ID取-将is_out设置为1;还-将is_out设置为0
func (lgc *Logics) TakeReturnByResID(useLog *UseLog) error {

	// 事务
	tx := lgc.db.Begin()
	if err := lgc.db.Model(&CabinetGrid{}).Where("in_res_id = ?", useLog.ResID).Update("is_out", useLog.Flag).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 修改使用状态，1-在库，2-借出
	status := 1
	if useLog.Flag == 1 {
		status = 2
	}
	if err := lgc.db.Model(&Sling{}).Where("id = ?", useLog.ResID).Update("use_status", status).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 记录借出日志
	if err := lgc.SaveTakeReturnLog(useLog); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// SaveTakeReturnLog 取还日志
func (lgc *Logics) SaveTakeReturnLog(useLog *UseLog) error {
	if useLog.Flag == 1 { // 借出，直接入库
		if err := lgc.db.Create(&useLog).Error; err != nil {
			return err
		}
	} else { // 归还，更新字段
		var log UseLog
		if err := lgc.db.Raw("SELECT * FROM t_res_use_log WHERE res_id = ? AND created_at = (SELECT MAX(created_at) FROM t_res_use_log WHERE res_id = ?)", useLog.ResID, useLog.ResID).
			Scan(&log).Error; err != nil {
			return err
		}
		if err := lgc.db.Model(&log).Updates(map[string]interface{}{"return_staff_id": useLog.ReturnStaffID, "return_staff_name": useLog.ReturnStaffName, "return_time": useLog.ReturnTime, "remark": useLog.Remark}).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetTakeReturnLog 取还日志
func (lgc *Logics) GetTakeReturnLog(param *UseLogQueryParam, pageIndex int, pageSize int) (*SearchResult, error) {

	logdb := lgc.db.Table("t_res_use_log").
		Select("id, res_name, take_staff_name, created_at, take_time, return_plan_time, return_staff_name, return_time, remark").
		// Select("t_res_use_log.*, t_res_sling.name AS res_name, t1.name AS take_staff_name, t2.name AS return_staff_name").
		// Joins("LEFT JOIN t_res_sling ON t_res_use_log.res_id = t_res_sling.id").
		// Joins("LEFT JOIN t_sys_staff AS t1 ON t_res_use_log.take_staff_id = t1.id").
		// Joins("LEFT JOIN t_sys_staff AS t2 ON t_res_use_log.return_staff_id = t2.id").
		Order("t_res_use_log.created_at desc")
	if param.ResName != "" {
		logdb = logdb.Where("t_res_use_log.res_name LIKE ?", "%"+param.ResName+"%")
	}
	if param.TakeStaff > 0 {
		logdb = logdb.Where("t_res_use_log.take_staff_id = ?", param.TakeStaff)
	}
	if param.ReturnStaff > 0 {
		logdb = logdb.Where("t_res_use_log.return_staff_id = ?", param.ReturnStaff)
	}
	if param.TakeStartTime != "" {
		logdb = logdb.Where("t_res_use_log.created_at >= ?", param.TakeStartTime)
	}
	if param.TakeEndTime != "" {
		logdb = logdb.Where("t_res_use_log.created_at <= ?", param.TakeEndTime)
	}
	if param.ReturnFlag == 1 { // 已归还
		logdb = logdb.Where("t_res_use_log.return_time IS NOT NULL")
	} else if param.ReturnFlag == 2 { // 未归还
		logdb = logdb.Where("t_res_use_log.return_time IS NULL")
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	logdb.Count(&rowCount)                                             //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var useLogs []UseLog
	if err := logdb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&useLogs).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &useLogs}, nil
}
