package models

import "gorm.io/gorm"

type ExchangeRate struct {
	gorm.Model
	// ID           uint      `gorm:"primarykey" json:"_id"`
	FromCurrency string  `json:"fromCurrency" binding:"required"`
	ToCurrency   string  `json:"toCurrency" binding:"required"`
	Rate         float64 `json:"rate" binding:"required"`
	// Date         time.Time `json:"date"`
}
