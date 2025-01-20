package gateway

import "gorm.io/gorm"

type gateway struct {
	db *gorm.DB
}

func NewGateway(db *gorm.DB) *gateway {
	return &gateway{db: db}
}
