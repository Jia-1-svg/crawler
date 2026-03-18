package model

import "gorm.io/gorm"

// 注册会员
type Member struct {
	gorm.Model
	Username string `gorm:"type:varchar(50)"`
	Phone    string `gorm:"type:varchar(50)"`
	Password string `gorm:"type:varchar(50)"`
	Status   int    `gorm:"type:int"`
}

/*
//电商商品
	type Goods struct {
		gorm.Model
		Title string  `gorm:"type:varchar(40)"`
		Price float64 `gorm:"type:decimal(10,2)"`
		Brand string  `gorm:"type:varchar(40)"`
		Status   int    `gorm:"type:int"`
		StockId int     `gorm:"type:int"`
	}
*/

// 积分
type Score struct {
	MemberId int    `gorm:"type:int"`
	Num      int    `gorm:"type:int"`
	Comment  string `gorm:"type:varchar(50)"`
}

// 库存
type GoodsStock struct {
	gorm.Model
	GoodsId int `gorm:"type:int"`
	Stock   int `gorm:"type:int"`
	Status  int `gorm:"type:int"`
}
