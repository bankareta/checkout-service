package entity

import "time"

// User is a struct that represents a user entity
type Etalase struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"column:username" json:"username"`
	ReferenceId  string     `gorm:"column:reference_id" json:"reference_id"`
	EtalaseName  string     `gorm:"column:etalase_name" json:"etalase_name"`
	PhotoEtalase string     `gorm:"column:photo_etalase" json:"photo_etalase,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"-"`
	TotalProduct int        `gorm:"-" json:"total_product"`
}

type EtalaseList struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"column:username" json:"username"`
	ReferenceId  string     `gorm:"column:reference_id" json:"reference_id"`
	EtalaseName  string     `gorm:"column:etalase_name" json:"etalase_name"`
	PhotoEtalase string     `gorm:"column:photo_etalase" json:"photo_etalase,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"-"`
	TotalProduct int        `gorm:"column:total_product" json:"total_product"`
}

type EtalaseUpdate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"column:username" json:"username"`
	ReferenceId  string    `gorm:"column:reference_id" json:"reference_id"`
	EtalaseName  string    `gorm:"column:etalase_name" json:"etalase_name"`
	PhotoEtalase string    `gorm:"column:photo_etalase" json:"photo_etalase"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type EtalaseDetail struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"column:username" json:"username"`
	ReferenceId  string     `gorm:"column:reference_id" json:"reference_id"`
	EtalaseName  string     `gorm:"column:etalase_name" json:"etalase_name"`
	PhotoEtalase string     `gorm:"column:photo_etalase" json:"photo_etalase,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (u *Etalase) TableName() string {
	return "etalases"
}

func (u *EtalaseUpdate) TableName() string {
	return "etalases"
}
