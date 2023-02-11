# develop handler layer

## PART1. 开发handler层

同样的,handler层需要实现`order/proto/order.pb.micro.go`中的`OrderHandler`接口.

`order/handler/order.go`:

```go
package handler

import (
	"context"
	"github.com/rayallen20/common"
	"github.com/rayallen20/order/domain/model"
	"github.com/rayallen20/order/domain/service"
	. "github.com/rayallen20/order/proto/order"
)

type Order struct {
	OrderDataService service.IOrderDataService
}

// GetOrderByID 根据订单ID查找订单信息
func (o *Order) GetOrderByID(ctx context.Context, request *OrderID, response *OrderInfo) error {
	order, err := o.OrderDataService.FindOrderByID(request.OrderId)
	if err != nil {
		return err
	}

	err = common.SwapTo(order, response)
	if err != nil {
		return err
	}

	return nil
}

// GetAllOrder 查询所有订单
func (o *Order) GetAllOrder(ctx context.Context, request *AllOrderRequest, response *AllOrder) error {
	orderAll, err := o.OrderDataService.FindAllOrder()
	if err != nil {
		return err
	}

	for _, v := range orderAll {
		order := &OrderInfo{}
		err = common.SwapTo(v, order)
		if err != nil {
			return err
		}
		response.OrderInfo = append(response.OrderInfo, order)
	}
	return nil
}

// CreateOrder 创建订单
func (o *Order) CreateOrder(ctx context.Context, request *OrderInfo, response *OrderID) error {
	orderAdd := &model.Order{}
	err := common.SwapTo(request, orderAdd)
	if err != nil {
		return err
	}

	orderID, err := o.OrderDataService.AddOrder(orderAdd)
	if err != nil {
		return err
	}

	response.OrderId = orderID
	return nil
}

// DeleteOrderByID 根据ID删除订单信息
func (o *Order) DeleteOrderByID(ctx context.Context, request *OrderID, response *Response) error {
	err := o.OrderDataService.DeleteOrder(request.OrderId)
	if err != nil {
		return err
	}
	response.Msg = "删除成功"
	return nil
}

// UpdateOrderPayStatus 更新订单支付状态
func (o *Order) UpdateOrderPayStatus(ctx context.Context, request *PayStatus, response *Response) error {
	err := o.OrderDataService.UpdatePayStatus(request.OrderId, request.PayStatus)
	if err != nil {
		return err
	}

	response.Msg = "支付状态更新成功"
	return nil
}

// UpdateOrderShipStatus 更新订单发货状态
func (o *Order) UpdateOrderShipStatus(ctx context.Context, request *ShipStatus, response *Response) error {
	err := o.OrderDataService.UpdateShipStatus(request.OrderId, request.ShipStatus)
	if err != nil {
		return err
	}

	response.Msg = "发货状态更新成功"
	return nil
}

// UpdateOrder 更新订单信息
func (o *Order) UpdateOrder(ctx context.Context, request *OrderInfo, response *Response) error {
	order := &model.Order{}
	err := common.SwapTo(request, order)
	if err != nil {
		return err
	}

	err = o.OrderDataService.UpdateOrder(order)
	if err != nil {
		return err
	}

	response.Msg = "订单更新成功"
	return nil
}
```