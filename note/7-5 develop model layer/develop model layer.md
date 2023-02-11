# develop model layer

## PART1. 定义订单ORM

`order/domain/model/order.go`:

```go
package model

// OrderDetail 订单详情ORM
type OrderDetail struct {
	// ID 订单详情id
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	// ProductID 订单详情所对应的商品信息ID
	ProductID int64 `json:"product_id"`
	// ProductNum 订单详情所对应的商品购买数量
	ProductNum int64 `json:"product_num"`
	// ProductSizeID 订单详情所对应商品的尺码信息ID
	ProductSizeID int64 `json:"product_size_id"`
	// ProductPrice 订单详情所对应商品的原价
	ProductPrice float64 `json:"product_price"`
	// OrderID 订单详情所属订单id
	OrderID int64 `json:"order_id"`
}
```

## PART2. 定义订单详情ORM

`order/domain/model/order_detail.go`:

```go
package model

import "time"

// Order 订单信息ORM
type Order struct {
	// ID 订单id
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	// OrderCode 用于满足订单幂等性 确保反复创建订单时 被创建的订单不会重复
	OrderCode string `gorm:"unique_index;not_null" json:"order_code"`
	// PayStatus 支付状态
	PayStatus int32 `json:"pay_status"`
	// ShipStatus 发货状态
	ShipStatus int32 `json:"ship_status"`
	// Price 订单价格
	Price float64 `json:"price"`
	// OrderDetail 订单详情集合
	OrderDetail []OrderDetail `gorm:"ForeignKey:OrderID" json:"order_detail"`
	// CreateAt 创建时间
	CreateAt time.Time
	// UpdateAt 修改时间
	UpdateAt time.Time
}
```