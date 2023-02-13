# develop handler layer

## PART1. model层开发

`payment/domain/model/payment.go`:

```go
package model

type Payment struct {
	// ID 主键自增ID
	ID int64 `gorm:"primary_key;not_null;auto_increment"`
	// PaymentName 支付单名称
	PaymentName string `json:"payment_name"`
	// PaymentSid 支付应用的secret
	PaymentSid string `json:"payment_sid"`
	// PaymentStatus 支付环境
	// true: live
	// false: sandbox
	PaymentStatus bool `json:"payment_status"`
	// PaymentImage 支付logo的url
	PaymentImage string `json:"payment_image"`
}
```

## PART2. handler层开发

`payment/handler/payment.go`:

```go
package handler

import (
	"context"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/model"
	"github.com/rayallen20/payment/domain/service"
	. "github.com/rayallen20/payment/proto/payment"
)

type Payment struct {
	PaymentDataService service.IPaymentDataService
}

// AddPayment 创建支付单
func (p *Payment) AddPayment(ctx context.Context, request *PaymentInfo, response *PaymentID) error {
	payment := &model.Payment{}
	err := common.SwapTo(request, payment)
	if err != nil {
		// 如有错误 记录日志 不再返回error
		common.Error(err)
	}

	paymentId, err := p.PaymentDataService.AddPayment(payment)
	if err != nil {
		common.Error(err)
	}

	response.PaymentId = paymentId
	return nil
}

// UpdatePayment 更新支付单信息
func (p *Payment) UpdatePayment(ctx context.Context, request *PaymentInfo, response *Response) error {
	payment := &model.Payment{}
	err := common.SwapTo(request, payment)
	if err != nil {
		common.Error(err)
	}

	err = p.PaymentDataService.UpdatePayment(payment)
	if err != nil {
		common.Error(err)
	}
	response.Msg = "更新成功"
	return nil
}

// DeletePayment 删除支付单信息
func (p *Payment) DeletePayment(ctx context.Context, request *PaymentID, response *Response) error {
	err := p.PaymentDataService.DeletePayment(request.PaymentId)
	if err != nil {
		common.Error(err)
	}
	response.Msg = "删除成功"
	return nil
}

// FindPaymentByID 根据ID查找支付单信息
func (p *Payment) FindPaymentByID(ctx context.Context, request *PaymentID, response *PaymentInfo) error {
	payment, err := p.PaymentDataService.FindPaymentByID(request.PaymentId)
	if err != nil {
		common.Error(err)
	}

	err = common.SwapTo(payment, response)
	if err != nil {
		common.Error(err)
	}

	return nil
}

// FindAllPayment 查找所有支付单信息
func (p *Payment) FindAllPayment(ctx context.Context, request *All, response *PaymentAll) error {
	allPayment, err := p.PaymentDataService.FindAllPayment()
	if err != nil {
		common.Error(err)
	}

	for _, payment := range allPayment {
		paymentInfo := &PaymentInfo{}
		err = common.SwapTo(payment, paymentInfo)
		if err != nil {
			common.Error(err)
		}
		response.PaymentInfo = append(response.PaymentInfo, paymentInfo)
	}

	return nil
}
```

