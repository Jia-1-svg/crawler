package model

/*
import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderSn  string  `gorm:"type:varchar(50);uniqueIndex"` //订单编号
	MemberID int     `gorm:"type:int"`                     //用户id
	Total    float64 `gorm:"type:decimal(10,2);default:0"` //总价
	PayType  int     `gorm:"type:int"`                     //支付方式
	Status   int     `gorm:"type:tinyint(1);default:0"`    //订单支付状态
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

func (o *Order) SaveOrder(db *gorm.DB) error {
	return db.Debug().Save(&o).Error
}

type OrderItem struct {
	gorm.Model
	OrderID  int `gorm:"type:int"`
	GoodsID  int `gorm:"type:int"`
	Quantity int `gorm:"type:int"`
}

// 会员
type Member struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex"`
	Password string `gorm:"type:varchar(255)"`
	Phone    string `gorm:"type:varchar(20)"`
	Nickname string `gorm:"type:varchar(50)"`
	Status   int    `gorm:"type:tinyint(1);default:1"` // 0-禁用 1-正常
}

// 积分
type Point struct {
	gorm.Model
	MemberID int    `gorm:"type:int"`
	OrderSn  string `gorm:"type:varchar(50)"`
	Points   int    `gorm:"type:int;default:0"` // 变动积分
	Balance  int    `gorm:"type:int;default:0"` // 当前余额
	Type     int    `gorm:"type:tinyint(1)"`    // 1-获取 2-消费
}

// 商品模型
type Goods struct {
	gorm.Model
	Title  string  `gorm:"type:varchar(200);not null"`
	Price  float64 `gorm:"type:decimal(10,2);default:0"`
	Points int     `gorm:"type:int;default:0"` // 赠送积分
	Image  string  `gorm:"type:varchar(500)"`
	Status int     `gorm:"type:tinyint(1);default:1"` // 0-下架 1-上架
	Stock  int     `gorm:"type:int;default:0"`        // 总库存
}

func (g *Goods) GoodsAdd(db *gorm.DB) error {
	return db.Debug().Create(&g).Error
}

func (g *Goods) DelGoodsById(db *gorm.DB, id int64) error {
	return db.Debug().Where("id =?", id).Delete(&g).Error
}

func (g *Goods) FindGoodsById(db *gorm.DB, id int64) error {
	return db.Debug().Where("id =?", id).First(&g).Error
}

func (g *Goods) GoodsSave(db *gorm.DB) error {
	return db.Debug().Save(&g).Error
}

// 库存
type Stock struct {
	gorm.Model
	GoodsID   int `gorm:"type:int;default:0"` // 库存数量
	Quantity  int `gorm:"type:int;default:0"` // 库存数量
	LockStock int `gorm:"type:int;default:0"` // 锁定库存
}

// 仓库
type Warehouse struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Phone   string `gorm:"type:varchar(20)"`
	Address string `gorm:"type:varchar(200)"`
	Status  int    `gorm:"type:tinyint(1);default:1"` // 0-禁用 1-正常
}

// 物流配送
type Logistics struct {
	gorm.Model
	OrderSn      string `gorm:"type:varchar(50);uniqueIndex"`
	Carrier      string `gorm:"type:varchar(50)"`          // 承运商
	Status       int    `gorm:"type:tinyint(1);default:0"` // 0-待发货 1-已发货 2-已签收
	ReceiverName string `gorm:"type:varchar(50)"`          // 收货人
	ReceiverAddr string `gorm:"type:varchar(200)"`         // 收货地址
}
*/
