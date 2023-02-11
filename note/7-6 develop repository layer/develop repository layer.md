# develop repository layer

## PART1. 实现repository层

`order/domain/repository/order_repository.go`:

```go
package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/rayallen20/order/domain/model"
)

type IOrderRepository interface {
	// InitTable 初始化数据库结构
	InitTable() error
	// FindOrderByID 根据ID查找订单
	FindOrderByID(int64) (*model.Order, error)
	// CreateOrder 创建订单
	CreateOrder(*model.Order) (int64, error)
	// DeleteOrderByID 根据ID删除订单
	DeleteOrderByID(int64) error
	// UpdateOrder 更新订单信息
	UpdateOrder(*model.Order) error
	// FindAll 查找所有订单
	FindAll() ([]model.Order, error)
	// UpdateShipStatus 更新发货状态
	UpdateShipStatus(int64, int32) error
	// UpdatePayStatus 更新支付状态
	UpdatePayStatus(int64, int32) error
}

// NewOrderRepository 创建orderRepository
func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{mysqlDb: db}
}

type OrderRepository struct {
	mysqlDb *gorm.DB
}

//InitTable 初始化数据库结构
func (o *OrderRepository) InitTable() error {
	return o.mysqlDb.CreateTable(&model.Order{}, &model.OrderDetail{}).Error
}

// FindOrderByID 根据ID查找订单
func (o *OrderRepository) FindOrderByID(orderID int64) (order *model.Order, err error) {
	order = &model.Order{}
	return order, o.mysqlDb.Preload("OrderDetail").First(order, orderID).Error
}

// CreateOrder 创建订单
func (o *OrderRepository) CreateOrder(order *model.Order) (int64, error) {
	return order.ID, o.mysqlDb.Create(order).Error
}

// DeleteOrderByID 根据ID删除订单
func (o *OrderRepository) DeleteOrderByID(orderID int64) error {
	tx := o.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 删除Order信息
	err := tx.Unscoped().Where("id = ?", orderID).Delete(&model.Order{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除OrderDetail信息
	err = tx.Unscoped().Where("order_id = ?", orderID).Delete(&model.OrderDetail{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateOrder 更新订单信息
func (o *OrderRepository) UpdateOrder(order *model.Order) error {
	return o.mysqlDb.Model(order).Update(order).Error
}

// FindAll 查找所有订单
func (o *OrderRepository) FindAll() (orderAll []model.Order, err error) {
	return orderAll, o.mysqlDb.Preload("OrderDetail").Find(&orderAll).Error
}

// UpdateShipStatus 更新发货状态
func (o *OrderRepository) UpdateShipStatus(orderID int64, shipStatus int32) error {
	db := o.mysqlDb.Model(&model.Order{}).Where("id = ?", orderID).UpdateColumn("ship_status", shipStatus)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}

// UpdatePayStatus 更新支付状态
func (o *OrderRepository) UpdatePayStatus(orderID int64, payStatus int32) error {
	db := o.mysqlDb.Model(&model.Order{}).Where("id = ?", orderID).UpdateColumn("pay_status", payStatus)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}
```