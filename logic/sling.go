package logic

import (
	"math"

	"zone.com/common"
)

// Sling 吊索具
type Sling struct {
	BaseModel
	RfID          string `json:"rfId" gorm:"size:64"`
	Name          string `json:"name" gorm:"size:64"` // 吊索具名称
	SlingType     uint   `json:"slingType"`
	MaxTonnage    uint   `json:"maxTonnage"`
	UseCount      int    `json:"useCount" gorm:"-"`
	UseStatus     uint   `json:"useStatus"`
	InspectStatus uint   `json:"inspectStatus"`
	PutTime       string `json:"putTime"`
	UsePermission string `json:"usePermission"`
	CabinetName   string `json:"cabinetName" gorm:"-"`
	CabinetID     uint   `json:"cabinetId" gorm:"-"`
	GridNo        uint   `json:"gridNo" gorm:"-"`
	IsOut         uint   `json:"isOut" gorm:"-"`
}

// TableName Sling
func (Sling) TableName() string {
	return "t_res_sling"
}

// ListSlings 查询吊索具
func (lgc *Logics) ListSlings(name string, slingType uint, maxTonnage uint, useStatus uint, inspectStatus uint, pageIndex int, pageSize int) (*SearchResult, error) {
	slingdb := lgc.db.Table("t_res_sling").
		Select("t_res_sling.*, t_res_cabinet.name AS cabinet_name, t_res_cabinet_grid.cabinet_id AS cabinet_id, t_res_cabinet_grid.grid_no AS grid_no, t_res_cabinet_grid.is_out AS is_out, t1.use_count").
		Joins("LEFT JOIN t_res_cabinet_grid ON t_res_cabinet_grid.in_res_id = t_res_sling.id").
		Joins("LEFT JOIN t_res_cabinet ON t_res_cabinet_grid.cabinet_id = t_res_cabinet.id").
		Joins("LEFT JOIN (SELECT t_res_use_log.res_id, COUNT(0) AS use_count FROM t_res_use_log GROUP BY t_res_use_log.res_id) t1 ON t1.res_id = t_res_sling.id").
		Where("t_res_sling.deleted_at IS NULL")
	if name != "" {
		slingdb = slingdb.Where("t_res_sling.name LIKE ?", "%"+name+"%")
	}
	if slingType > 0 {
		slingdb = slingdb.Where("t_res_sling.sling_type = ?", slingType)
	}
	if maxTonnage > 0 {
		slingdb = slingdb.Where("t_res_sling.max_tonnage = ?", maxTonnage)
	}
	if useStatus > 0 {
		slingdb = slingdb.Where("t_res_sling.use_status = ?", useStatus)
	}
	if inspectStatus > 0 {
		slingdb = slingdb.Where("t_res_sling.inspect_status = ?", inspectStatus)
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	slingdb.Count(&rowCount)                                           //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var slings []Sling
	if err := slingdb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&slings).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &slings}, nil
}

// AddSling 添加吊索具
func (lgc *Logics) AddSling(sling *Sling) error {
	// 吊索具名字不能为空
	if sling.Name == "" {
		return common.ErrSlingNameIsNull
	}
	// 吊索具RFID不能为空
	if sling.RfID == "" {
		return common.ErrSlingRfIDIsNull
	}
	// 存放位置为空
	// if sling.CabinetID == 0 || sling.GridNo == 0 {
	// 	return  utils.ErrSlingCabinetIsNull
	// }
	// 名字或RFID重复
	sling0, _ := lgc.QuerySlingByName(sling.Name, sling.RfID)
	if sling0 != nil {
		return common.ErrSlingAlreadyExists
	}
	if err := lgc.db.Create(&sling).Error; err != nil {
		return err
	}
	// 保存位置数据
	// sling1, _ := s.QuerySlingByName(sling.Name, sling.RfID)
	// _, err1 := s.Store(sling.CabinetID, sling.GridNo, sling1.ID)
	// if err1 != nil { // 失败删除
	// 	lgc.db.Unscoped().Where("id = ?", sling1.ID).Delete(&Sling{})
	// 	return "", err1
	// }
	return nil
}

// QuerySlingByName 查询吊索具
func (lgc *Logics) QuerySlingByName(name, rfID string) (*Sling, error) {
	var sling Sling
	if err := lgc.db.Where("name = ? OR rf_id = ?", name, rfID).First(&sling).Error; err != nil {
		return nil, err
	}

	return &sling, nil
}

// UpdateSling 修改吊索具
func (lgc *Logics) UpdateSling(sling *Sling) error {
	// 吊索具RFID不能为空
	if sling.RfID == "" {
		return common.ErrSlingRfIDIsNull
	}
	// 吊索具名字不能为空
	if sling.Name == "" {
		return common.ErrSlingNameIsNull
	}
	// 存放位置为空
	// if sling.CabinetID == 0 || sling.GridNo == 0 {
	// 	return  utils.ErrSlingCabinetIsNull
	// }
	// 名字或RFID重复
	sling0, _ := lgc.QuerySlingByName(sling.Name, sling.RfID)
	if sling0 != nil && sling0.ID != sling.ID {
		return common.ErrSlingAlreadyExists
	}
	// 事务
	tx := lgc.db.Begin()
	if err := tx.Save(&sling).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 保存位置数据
	// _, err1 := s.Store(sling.CabinetID, sling.GridNo, sling.ID)
	// if err1 != nil {
	// 	tx.Rollback()
	// 	return  err1
	// }
	tx.Commit()
	return nil
}

// DeleteSling 删除吊索具
func (lgc *Logics) DeleteSling(id uint) error {
	// 事务
	tx := lgc.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&Sling{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除占用的箱格
	if err := tx.Where("in_res_id = ?", id).Delete(&CabinetGrid{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
