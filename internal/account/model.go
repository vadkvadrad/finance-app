package account

import "gorm.io/gorm"

const (
	CurrencyRub = "RUB"
)

type Account struct {
	gorm.Model
	UserID   int     `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}