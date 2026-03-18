package request

type GoodsAdd struct {
	Title string  `form:"title" json:"title" xml:"title"  binding:"required"`
	Price float64 `form:"price" json:"price" xml:"price"  binding:"required"`
	Stock int     `form:"stock" json:"stock" xml:"stock" binding:"required"`
}
type GoodsDel struct {
	GoodsId int `form:"goodsId" json:"goodsId" xml:"goodsId"  binding:"required"`
}
type GoodsUpdate struct {
	GoodsId int     `form:"goodsId" json:"goodsId" xml:"goodsId"  binding:"required"`
	Title   string  `form:"title" json:"title" xml:"title"  binding:"required"`
	Price   float64 `form:"price" json:"price" xml:"price"  binding:"required"`
	Stock   int     `form:"stock" json:"stock" xml:"stock" binding:"required"`
}
type GoodsDetail struct {
	GoodsId int `form:"goodsId" json:"goodsId" xml:"goodsId"  binding:"required"`
}
type OrderAdd struct {
	UserId    int            `json:"userId" binding:"required"`
	PayType   int            `json:"payType"   binding:"required"`
	OrderItem []OrderItemAdd `json:"orderItem"` // 将 OrderItemAdd 作为具名结构体字段
}
type OrderItemAdd struct {
	GoodsId  int `json:"goodsId"   binding:"required"`
	Quantity int `json:"quantity"   binding:"required"`
}
