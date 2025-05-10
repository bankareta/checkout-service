package entity

// User is a struct that represents a user entity
type EtalaseProduct struct {
	EtalaseID uint `gorm:"column:id_etalase"`
	ProductID uint `gorm:"column:id_product"`
}

func (u *EtalaseProduct) TableName() string {
	return "etalase_product"
}
