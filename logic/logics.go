package logic

import (
	"gorm.io/gorm"
)

// Logics framwork need
type Logics struct {
	db *gorm.DB
}

func NewLogics(db *gorm.DB) *Logics {
	return &Logics{db: db}
}
