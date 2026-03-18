package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nickname string `gorm:"type:varchar(40)"`
	Password string `gorm:"type:varchar(40)"`
}

func (u *User) UserAdd(db *gorm.DB) error {
	return db.Debug().Create(&u).Error
}

type Goods struct {
	gorm.Model
	Title string  `gorm:"type:varchar(40)"`
	Price float64 `gorm:"type:decimal(10,2)"`
	Stock int     `gorm:"type:int"`
}

func (g *Goods) GoodsAdd(db *gorm.DB) error {
	return db.Debug().Create(&g).Error
}

func (g *Goods) GoodsDel(db *gorm.DB, id int64) error {
	return db.Debug().Where("id =?", id).Delete(&g).Error
}

func (g *Goods) FindGoodsById(db *gorm.DB, id int64) error {
	return db.Debug().Where("id =?", id).First(&g).Error
}

func (g *Goods) UpdateGoods(db *gorm.DB) error {
	return db.Debug().Save(&g).Error
}
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type Order struct {
	gorm.Model
	UserId  int     `gorm:"type:int"`
	PayType int     `gorm:"type:int"`
	OrderSn string  `gorm:"type:varchar(50)"`
	Status  int     `gorm:"type:int"`
	Total   float64 `gorm:"type:decimal(10,2)"`
}

func (o *Order) OrderAdd(db *gorm.DB) error {
	return db.Debug().Create(&o).Error
}

func (o *Order) OrderItemAdd(db *gorm.DB, item []*OrderItem) error {
	return db.Debug().Create(&item).Error
}

func (o *Order) FindOrderByOrderSn(db *gorm.DB, sn string) error {
	return db.Debug().Where("order_sn =?", sn).First(&o).Error
}

func (o *Order) UpdateOrderStatus(db *gorm.DB) error {
	return db.Debug().Save(&o).Error
}

type OrderItem struct {
	gorm.Model
	OrderId  int `gorm:"type:int"`
	GoodsId  int `gorm:"type:int"`
	Quantity int `gorm:"type:int"`
}

func (i *OrderItem) ItemAdd(db *gorm.DB) error {
	return db.Debug().Create(&i).Error
}
