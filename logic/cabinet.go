package logic

import (
	"math"

	"zone.com/common"
)

// Cabinet 智能柜
type Cabinet struct {
	BaseModel
	Name        string `json:"name" gorm:"size:64"` // 智能柜名称
	GridCount   uint   `json:"gridCount"`
	Location    string `json:"location"`
	UsedCount   uint   `json:"usedCount" gorm:"-"`
	UnUsedCount uint   `json:"unUsedCount" gorm:"-"`
	Status      int16  `json:"status"` // 状态：0-正常
	Remark      string `json:"remark"` // 说明
}

// TableName Cabinet
func (Cabinet) TableName() string {
	return "t_res_cabinet"
}

// CabinetGrid 智能柜箱格
type CabinetGrid struct {
	BaseModel
	GridNo    uint `json:"gridNo"`    // 箱格编号
	CabinetID uint `json:"cabinetId"` // 智能柜ID
	InResID   uint `json:"inResId"`   // 存放的资产ID，空为0
	IsOut     uint `json:"isOut"`     // 是否借出
}

// TableName CabinetGrid
func (CabinetGrid) TableName() string {
	return "t_res_cabinet_grid"
}

// ListCabinets 查询智能柜
func (lgc *Logics) ListCabinets(name string, pageIndex int, pageSize int) (*SearchResult, error) {

	cabinetdb := lgc.db.Table("t_res_cabinet").
		Select("t_res_cabinet.*, COALESCE(t1.used_count, 0) AS used_count, COALESCE(t_res_cabinet.grid_count - t1.used_count, t_res_cabinet.grid_count) AS un_used_count").
		Joins("LEFT JOIN (SELECT t_res_cabinet_grid.cabinet_id, COUNT(0) AS used_count FROM t_res_cabinet_grid WHERE t_res_cabinet_grid.in_res_id > 0 AND t_res_cabinet_grid.deleted_at IS NULL GROUP BY cabinet_id) t1 ON t1.cabinet_id = t_res_cabinet.id").
		Where("t_res_cabinet.deleted_at IS NULL")
	if name != "" {
		cabinetdb = cabinetdb.Where("t_res_cabinet.name LIKE ?", "%"+name+"%")
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	var rowCount int64
	cabinetdb.Count(&rowCount)                                         //总行数
	pageCount := int(math.Ceil(float64(rowCount) / float64(pageSize))) // 总页数

	var cabinets []Cabinet
	if err := cabinetdb.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&cabinets).Error; err != nil {
		return nil, err
	}

	return &SearchResult{Total: rowCount, PageIndex: pageIndex, PageSize: pageSize, PageCount: pageCount, List: &cabinets}, nil
}

// AddCabinet 添加智能柜
func (lgc *Logics) AddCabinet(cabinet *Cabinet) error {
	// 智能柜名字不能为空
	if cabinet.Name == "" {
		return common.ErrCabinetNameIsNull
	}

	// 智能柜箱格不能是0
	if cabinet.GridCount == 0 {
		return common.ErrCabinetGridIsZero
	}
	// 名字重复
	cabinet0, _ := lgc.QueryCabinetByName(cabinet.Name)
	if cabinet0 != nil {
		return common.ErrCabinetAlreadyExists
	}
	if err := lgc.db.Create(&cabinet).Error; err != nil {
		return err
	}
	return nil
}

// QueryCabinetByName 查询智能柜
func (lgc *Logics) QueryCabinetByName(name string) (*Cabinet, error) {

	var cabinet Cabinet
	if err := lgc.db.Where("name = ?", name).First(&cabinet).Error; err != nil {
		return nil, err
	}

	return &cabinet, nil
}

// QueryCabinetByID 查询智能柜
func (lgc *Logics) QueryCabinetByID(id uint) (*Cabinet, error) {

	var cabinet Cabinet
	if err := lgc.db.Where("id = ?", id).First(&cabinet).Error; err != nil {
		return nil, err
	}

	return &cabinet, nil
}

// UpdateCabinet 修改智能柜
func (lgc *Logics) UpdateCabinet(cabinet *Cabinet) error {

	// 智能柜名字不能为空
	if cabinet.Name == "" {
		return common.ErrCabinetNameIsNull
	}

	// 智能柜箱格不能是0
	if cabinet.GridCount == 0 {
		return common.ErrCabinetGridIsZero
	}
	// 名字重复
	cabinet0, _ := lgc.QueryCabinetByName(cabinet.Name)
	if cabinet0 != nil && cabinet0.ID != cabinet.ID {
		return common.ErrCabinetAlreadyExists
	}
	if err := lgc.db.Save(&cabinet).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCabinet 删除智能柜
func (lgc *Logics) DeleteCabinet(id uint) error {
	// 事务
	tx := lgc.db.Begin()
	// 删除智能柜
	if err := tx.Where("id = ?", id).Delete(&Cabinet{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除箱格
	if err := tx.Where("cabinet_id = ?", id).Delete(&CabinetGrid{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// ListCabinetGrids 查询箱格列表
func (lgc *Logics) ListCabinetGrids(cabinetID uint) (*SearchResult, error) {

	// 智能柜
	cabinet, err := lgc.QueryCabinetByID(cabinetID)
	if err != nil || cabinet == nil {
		return nil, common.ErrNotFound
	}

	// 在用
	griddb := lgc.db.Model(&CabinetGrid{})
	if cabinetID > 0 {
		griddb = griddb.Where("cabinet_id = ?", cabinetID)
	}

	var grids []CabinetGrid
	if err := griddb.Find(&grids).Error; err != nil {
		return nil, err
	}

	result := make([]CabinetGrid, cabinet.GridCount)
	for i := 0; i < int(cabinet.GridCount); i++ {
		flag := false
		for _, v := range grids {
			if int(v.GridNo) == i+1 { // 已使用
				data := &CabinetGrid{GridNo: uint(i + 1), CabinetID: cabinetID, InResID: v.InResID}
				result[i] = *data
				flag = true
				break
			}
		}
		if !flag { // 未使用
			data := &CabinetGrid{GridNo: uint(i + 1), CabinetID: cabinetID, InResID: 0}
			result[i] = *data
		}
	}

	return &SearchResult{Total: int64(len(result)), PageIndex: 0, PageSize: 0, PageCount: 0, List: &result}, nil
}
