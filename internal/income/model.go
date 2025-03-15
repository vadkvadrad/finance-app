package income

import "gorm.io/gorm"

type Income struct {
	gorm.Model
	UserId uint `gorm:"index"`
	Amount float64
}
