# devlop model layer

## PART1. 定义ORM

`cart/domain/model/cart.go`代码如下:

```go
package model

// Cart 购物车ORM
type Cart struct {
	// ID 主键自增id
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// ProductID 商品id
	ProductID int64 `gorm:"not_null" json:"product_id"`

	// Num 商品数量
	Num int64 `gorm:"not_null" json:"num"`

	// SizeID 尺码id
	SizeID int64 `gorm:"not_null" json:"size_id"`

	// UserID 用户id
	UserID int64 `gorm:"not_null" json:"user_id"`
}
```