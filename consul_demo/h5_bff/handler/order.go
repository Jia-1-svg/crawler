package handler

import (
	"net/http"
	"strconv"

	"consul_demo/h5_bff/client"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// OrderHandler 订单 HTTP 处理器
type OrderHandler struct {
	orderClient *client.OrderClient
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderClient *client.OrderClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID     int64   `json:"user_id" binding:"required"`
	TotalPrice float64 `json:"total_price" binding:"required"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID int64  `json:"order_id"`
	OrderNo string `json:"order_no"`
}

// CreateOrderHandler 创建订单
func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
			Data:    nil,
		})
		return
	}

	orderID, orderNo, err := h.orderClient.CreateOrder(c, req.UserID, req.TotalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: CreateOrderResponse{
			OrderID: orderID,
			OrderNo: orderNo,
		},
	})
}

// GetOrderHandler 获取订单详情
func (h *OrderHandler) GetOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid order ID: " + err.Error(),
			Data:    nil,
		})
		return
	}

	order, err := h.orderClient.GetOrder(c, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    order,
	})
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	TotalPrice float64 `json:"total_price" binding:"required"`
	Status     int32   `json:"status" binding:"required"`
}

// UpdateOrderHandler 更新订单
func (h *OrderHandler) UpdateOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid order ID: " + err.Error(),
			Data:    nil,
		})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
			Data:    nil,
		})
		return
	}

	if err := h.orderClient.UpdateOrder(c, orderID, req.TotalPrice, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}

// DeleteOrderHandler 删除订单
func (h *OrderHandler) DeleteOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid order ID: " + err.Error(),
			Data:    nil,
		})
		return
	}

	if err := h.orderClient.DeleteOrder(c, orderID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}

// ListOrdersResponse 订单列表响应
type ListOrdersResponse struct {
	Orders []*Order `json:"orders"`
	Total  int32    `json:"total"`
}

// Order 订单响应结构
type Order struct {
	ID         int64   `json:"id"`
	OrderNo    string  `json:"order_no"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     int32   `json:"status"`
	CreateTime string  `json:"create_time"`
	UpdateTime string  `json:"update_time"`
}

// ListOrdersHandler 获取订单列表
func (h *OrderHandler) ListOrdersHandler(c *gin.Context) {
	// 解析查询参数
	userIDStr := c.Query("user_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "user_id is required",
			Data:    nil,
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid user_id: " + err.Error(),
			Data:    nil,
		})
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	orders, total, err := h.orderClient.ListOrders(c, userID, int32(page), int32(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// 转换为响应格式
	var orderList []*Order
	for _, o := range orders {
		orderList = append(orderList, &Order{
			ID:         o.Id,
			OrderNo:    o.OrderNo,
			UserID:     o.UserId,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreateTime: o.CreateTime,
			UpdateTime: o.UpdateTime,
		})
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: ListOrdersResponse{
			Orders: orderList,
			Total:  total,
		},
	})
}
