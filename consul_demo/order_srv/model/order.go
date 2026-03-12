package model

import "time"

type Order struct {
	ID         int64   `gorm:"column:id;primaryKey" json:"id"`
	OrderNo    string  `gorm:"column:order_no" json:"order_no"`
	UserID     int64   `gorm:"column:user_id" json:"user_id"`
	TotalPrice float64 `gorm:"column:total_price" json:"total_price"`
	Status     int32   `gorm:"column:status" json:"status"`
	CreateTime string  `gorm:"column:create_time" json:"create_time"`
	UpdateTime string  `gorm:"column:update_time" json:"update_time"`
}

const STATUS_CREATE = 1

func (Order) TableName() string {
	return "eb_store_order"
}

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
