# develop service layer

## PART1. 开发service层

`order/domain/service/order_data_service.go`:

```go
package service

import (
	"github.com/rayallen20/order/domain/model"
	"github.com/rayallen20/order/domain/repository"
)

type IOrderDataService interface {
	// AddOrder 新增订单
	AddOrder(*model.Order) (int64, error)
	// DeleteOrder 删除订单
	DeleteOrder(int64) error
	// UpdateOrder 更新订单
	UpdateOrder(*model.Order) error
	// FindOrderByID 根据ID查找订单
	FindOrderByID(int64) (*model.Order, error)
	// FindAllOrder 查找所有订单
	FindAllOrder() ([]model.Order, error)
	// UpdateShipStatus 更新发货状态
	UpdateShipStatus(int64, int32) error
	// UpdatePayStatus 更新支付状态
	UpdatePayStatus(int64, int32) error
}

// NewOrderDataService 创建Service层对象
func NewOrderDataService(orderRepository repository.IOrderRepository) IOrderDataService {
	return &OrderDataService{orderRepository}
}

type OrderDataService struct {
	OrderRepository repository.IOrderRepository
}

// AddOrder 新增订单
func (o *OrderDataService) AddOrder(order *model.Order) (int64, error) {
	return o.OrderRepository.CreateOrder(order)
}

// DeleteOrder 删除订单
func (o *OrderDataService) DeleteOrder(orderID int64) error {
	return o.OrderRepository.DeleteOrderByID(orderID)
}

// UpdateOrder 更新订单
func (o *OrderDataService) UpdateOrder(order *model.Order) error {
	return o.OrderRepository.UpdateOrder(order)
}

// FindOrderByID 根据ID查找订单
func (o *OrderDataService) FindOrderByID(orderID int64) (*model.Order, error) {
	return o.OrderRepository.FindOrderByID(orderID)
}

// FindAllOrder 查找所有订单
func (o *OrderDataService) FindAllOrder() ([]model.Order, error) {
	return o.OrderRepository.FindAll()
}

// UpdateShipStatus 更新发货状态
func (o *OrderDataService) UpdateShipStatus(orderID int64, shipStatus int32) error {
	return o.OrderRepository.UpdateShipStatus(orderID, shipStatus)
}

// UpdatePayStatus 更新支付状态
func (o *OrderDataService) UpdatePayStatus(orderID int64, payStatus int32) error {
	return o.OrderRepository.UpdatePayStatus(orderID, payStatus)
}
```