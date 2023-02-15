# develop handler of paymentApi--second half

## PART1. 编写退款信息

`paymentApi/handler/paymentApi.go`:

```go
package handler

import (
	"context"
	"errors"
	"github.com/plutov/paypal/v3"
	"github.com/rayallen20/common"
	go_micro_service_payment "github.com/rayallen20/payment/proto/payment"
	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
	"strconv"
)

type PaymentApi struct {
	PaymentService go_micro_service_payment.PaymentService
}

var (
	// ClientID PayPal的sandbox app的id
	ClientID string = "AVKXYfKLOTMmueeQeQeDH5skLWbEgelOJYMYUOtkI97wm6n178zc6SUUpWScMC4y4cilrGrD8mxWSaMh"
)

// PayPalRefund PaymentApi.PayPalRefund 通过API向外暴露为/paymentApi/payPalRefund，接收http请求
// 即：/paymentApi/payPalRefund请求会调用go.micro.api.paymentApi 服务的PaymentApi.PayPalRefund方法
// 退款业务逻辑
func (e *PaymentApi) PayPalRefund(ctx context.Context, req *paymentApi.Request, rsp *paymentApi.Response) error {
	//验证payment支付通道是否赋值
	if err := isOK("payment_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款号
	if err := isOK("refund_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款金额
	if err := isOK("money", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	// 获取paymentID
	paymentID, err := strconv.ParseInt(req.Get["payment_id"].Values[0], 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取支付通道信息
	paymentInfo, err := e.PaymentService.FindPaymentByID(ctx, &go_micro_service_payment.PaymentID{PaymentId: paymentID})
	if err != nil {
		common.Error(err)
		return err
	}

	// 确认支付模式
	status := paypal.APIBaseSandBox
	// 若DB中记录为生产环境 则将status置为生产环境
	if paymentInfo.PaymentStatus {
		status = paypal.APIBaseLive
	}

	// 创建退款信息
	payout := paypal.Payout{
		// 发送电子邮件信息
		SenderBatchHeader: &paypal.SenderBatchHeader{
			// 电子邮件标题
			EmailSubject: req.Get["refund_id"].Values[0] + "imooc 提醒你收款!",
			// 电子邮件正文
			EmailMessage: req.Get["refund_id"].Values[0] + "您有一个收款信息!",
			// 每笔转账的唯一ID
			SenderBatchID: req.Get["refund_id"].Values[0],
		},
		// 付款信息
		Items: []paypal.PayoutItem{
			{
				// 接收提醒方式
				RecipientType: "EMAIL",
				// 退款的接收者 即personal sandbox account
				Receiver: "sb-q1kid25051837@personal.example.com",
				// 金额
				Amount: &paypal.AmountPayout{
					// 币种
					Currency: "USD",
					// 金额
					Value: req.Get["money"].Values[0],
				},
				// 备注
				Note: req.Get["refund_id"].Values[0],
				// 唯一ID
				SenderItemID: req.Get["refund_id"].Values[0],
			},
		},
	}
}

// isOK 根据给定的key名 判断该key在请求中是否存在
func isOK(key string, req *paymentApi.Request) error {
	if _, ok := req.Get[key]; !ok {
		err := errors.New(key + " 参数异常")
		common.Error(err)
		return err
	}
	return nil
}
```

## PART2. 创建支付客户端并获取token

### 2.1 创建支付客户端

`paymentApi/handler/paymentApi.go`:

```go
package handler

import (
	"context"
	"errors"
	"github.com/plutov/paypal/v3"
	"github.com/rayallen20/common"
	go_micro_service_payment "github.com/rayallen20/payment/proto/payment"
	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
	"strconv"
)

type PaymentApi struct {
	PaymentService go_micro_service_payment.PaymentService
}

var (
	// ClientID PayPal的sandbox app的id
	ClientID string = "AVKXYfKLOTMmueeQeQeDH5skLWbEgelOJYMYUOtkI97wm6n178zc6SUUpWScMC4y4cilrGrD8mxWSaMh"
)

// PayPalRefund PaymentApi.PayPalRefund 通过API向外暴露为/paymentApi/payPalRefund，接收http请求
// 即：/paymentApi/payPalRefund请求会调用go.micro.api.paymentApi 服务的PaymentApi.PayPalRefund方法
// 退款业务逻辑
func (e *PaymentApi) PayPalRefund(ctx context.Context, req *paymentApi.Request, rsp *paymentApi.Response) error {
	//验证payment支付通道是否赋值
	if err := isOK("payment_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款号
	if err := isOK("refund_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款金额
	if err := isOK("money", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	// 获取paymentID
	paymentID, err := strconv.ParseInt(req.Get["payment_id"].Values[0], 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取支付通道信息
	paymentInfo, err := e.PaymentService.FindPaymentByID(ctx, &go_micro_service_payment.PaymentID{PaymentId: paymentID})
	if err != nil {
		common.Error(err)
		return err
	}

	// 确认支付模式
	status := paypal.APIBaseSandBox
	// 若DB中记录为生产环境 则将status置为生产环境
	if paymentInfo.PaymentStatus {
		status = paypal.APIBaseLive
	}

	// 创建退款信息
	payout := paypal.Payout{
		// 发送电子邮件信息
		SenderBatchHeader: &paypal.SenderBatchHeader{
			// 电子邮件标题
			EmailSubject: req.Get["refund_id"].Values[0] + "imooc 提醒你收款!",
			// 电子邮件正文
			EmailMessage: req.Get["refund_id"].Values[0] + "您有一个收款信息!",
			// 每笔转账的唯一ID
			SenderBatchID: req.Get["refund_id"].Values[0],
		},
		// 付款信息
		Items: []paypal.PayoutItem{
			{
				// 接收提醒方式
				RecipientType: "EMAIL",
				// 退款的接收者 即personal sandbox account
				Receiver: "sb-q1kid25051837@personal.example.com",
				// 金额
				Amount: &paypal.AmountPayout{
					// 币种
					Currency: "USD",
					// 金额
					Value: req.Get["money"].Values[0],
				},
				// 备注
				Note: req.Get["refund_id"].Values[0],
				// 唯一ID
				SenderItemID: req.Get["refund_id"].Values[0],
			},
		},
	}

	// 创建支付客户端
	payPalClient, err := paypal.NewClient(ClientID, paymentInfo.PaymentSid, status)
	if err != nil {
		common.Error(err)
		return err
	}
}

// isOK 根据给定的key名 判断该key在请求中是否存在
func isOK(key string, req *paymentApi.Request) error {
	if _, ok := req.Get[key]; !ok {
		err := errors.New(key + " 参数异常")
		common.Error(err)
		return err
	}
	return nil
}
```

### 2.2 获取token

`paymentApi/handler/paymentApi.go`:

```go
package handler

import (
	"context"
	"errors"
	"github.com/plutov/paypal/v3"
	"github.com/rayallen20/common"
	go_micro_service_payment "github.com/rayallen20/payment/proto/payment"
	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
	"strconv"
)

type PaymentApi struct {
	PaymentService go_micro_service_payment.PaymentService
}

var (
	// ClientID PayPal的sandbox app的id
	ClientID string = "AVKXYfKLOTMmueeQeQeDH5skLWbEgelOJYMYUOtkI97wm6n178zc6SUUpWScMC4y4cilrGrD8mxWSaMh"
)

// PayPalRefund PaymentApi.PayPalRefund 通过API向外暴露为/paymentApi/payPalRefund，接收http请求
// 即：/paymentApi/payPalRefund请求会调用go.micro.api.paymentApi 服务的PaymentApi.PayPalRefund方法
// 退款业务逻辑
func (e *PaymentApi) PayPalRefund(ctx context.Context, req *paymentApi.Request, rsp *paymentApi.Response) error {
	//验证payment支付通道是否赋值
	if err := isOK("payment_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款号
	if err := isOK("refund_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款金额
	if err := isOK("money", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	// 获取paymentID
	paymentID, err := strconv.ParseInt(req.Get["payment_id"].Values[0], 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取支付通道信息
	paymentInfo, err := e.PaymentService.FindPaymentByID(ctx, &go_micro_service_payment.PaymentID{PaymentId: paymentID})
	if err != nil {
		common.Error(err)
		return err
	}

	// 确认支付模式
	status := paypal.APIBaseSandBox
	// 若DB中记录为生产环境 则将status置为生产环境
	if paymentInfo.PaymentStatus {
		status = paypal.APIBaseLive
	}

	// 创建退款信息
	payout := paypal.Payout{
		// 发送电子邮件信息
		SenderBatchHeader: &paypal.SenderBatchHeader{
			// 电子邮件标题
			EmailSubject: req.Get["refund_id"].Values[0] + "imooc 提醒你收款!",
			// 电子邮件正文
			EmailMessage: req.Get["refund_id"].Values[0] + "您有一个收款信息!",
			// 每笔转账的唯一ID
			SenderBatchID: req.Get["refund_id"].Values[0],
		},
		// 付款信息
		Items: []paypal.PayoutItem{
			{
				// 接收提醒方式
				RecipientType: "EMAIL",
				// 退款的接收者 即personal sandbox account
				Receiver: "sb-q1kid25051837@personal.example.com",
				// 金额
				Amount: &paypal.AmountPayout{
					// 币种
					Currency: "USD",
					// 金额
					Value: req.Get["money"].Values[0],
				},
				// 备注
				Note: req.Get["refund_id"].Values[0],
				// 唯一ID
				SenderItemID: req.Get["refund_id"].Values[0],
			},
		},
	}

	// 创建支付客户端
	payPalClient, err := paypal.NewClient(ClientID, paymentInfo.PaymentSid, status)
	if err != nil {
		common.Error(err)
		return err
	}
	
	// 获取token
	_, err = payPalClient.GetAccessToken()
	if err != nil {
		common.Error(err)
		return err
	}
}

// isOK 根据给定的key名 判断该key在请求中是否存在
func isOK(key string, req *paymentApi.Request) error {
	if _, ok := req.Get[key]; !ok {
		err := errors.New(key + " 参数异常")
		common.Error(err)
		return err
	}
	return nil
}
```

## PART3. 发送转账请求给paypal并填充响应

`paymentApi/handler/paymentApi.go`:

```go
package handler

import (
	"context"
	"errors"
	"github.com/plutov/paypal/v3"
	"github.com/rayallen20/common"
	go_micro_service_payment "github.com/rayallen20/payment/proto/payment"
	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
	"strconv"
)

type PaymentApi struct {
	PaymentService go_micro_service_payment.PaymentService
}

var (
	// ClientID PayPal的sandbox app的id
	ClientID string = "AVKXYfKLOTMmueeQeQeDH5skLWbEgelOJYMYUOtkI97wm6n178zc6SUUpWScMC4y4cilrGrD8mxWSaMh"
)

// PayPalRefund PaymentApi.PayPalRefund 通过API向外暴露为/paymentApi/payPalRefund，接收http请求
// 即：/paymentApi/payPalRefund请求会调用go.micro.api.paymentApi 服务的PaymentApi.PayPalRefund方法
// 退款业务逻辑
func (e *PaymentApi) PayPalRefund(ctx context.Context, req *paymentApi.Request, rsp *paymentApi.Response) error {
	//验证payment支付通道是否赋值
	if err := isOK("payment_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款号
	if err := isOK("refund_id", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	//验证退款金额
	if err := isOK("money", req); err != nil {
		rsp.StatusCode = 500
		return err
	}

	// 获取paymentID
	paymentID, err := strconv.ParseInt(req.Get["payment_id"].Values[0], 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取支付通道信息
	paymentInfo, err := e.PaymentService.FindPaymentByID(ctx, &go_micro_service_payment.PaymentID{PaymentId: paymentID})
	if err != nil {
		common.Error(err)
		return err
	}

	// 确认支付模式
	status := paypal.APIBaseSandBox
	// 若DB中记录为生产环境 则将status置为生产环境
	if paymentInfo.PaymentStatus {
		status = paypal.APIBaseLive
	}

	// 创建退款信息
	payout := paypal.Payout{
		// 发送电子邮件信息
		SenderBatchHeader: &paypal.SenderBatchHeader{
			// 电子邮件标题
			EmailSubject: req.Get["refund_id"].Values[0] + "imooc 提醒你收款!",
			// 电子邮件正文
			EmailMessage: req.Get["refund_id"].Values[0] + "您有一个收款信息!",
			// 每笔转账的唯一ID
			SenderBatchID: req.Get["refund_id"].Values[0],
		},
		// 付款信息
		Items: []paypal.PayoutItem{
			{
				// 接收提醒方式
				RecipientType: "EMAIL",
				// 退款的接收者 即personal sandbox account
				Receiver: "sb-q1kid25051837@personal.example.com",
				// 金额
				Amount: &paypal.AmountPayout{
					// 币种
					Currency: "USD",
					// 金额
					Value: req.Get["money"].Values[0],
				},
				// 备注
				Note: req.Get["refund_id"].Values[0],
				// 唯一ID
				SenderItemID: req.Get["refund_id"].Values[0],
			},
		},
	}

	// 创建支付客户端
	payPalClient, err := paypal.NewClient(ClientID, paymentInfo.PaymentSid, status)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取token
	_, err = payPalClient.GetAccessToken()
	if err != nil {
		common.Error(err)
		return err
	}

	// 发送转账请求给paypal
	paymentResult, err := payPalClient.CreateSinglePayout(payout)
	if err != nil {
		common.Error(err)
		return err
	}
	common.Info(paymentResult)
	rsp.Body = req.Get["refund_id"].Values[0] + "支付成功!"
	return nil
}

// isOK 根据给定的key名 判断该key在请求中是否存在
func isOK(key string, req *paymentApi.Request) error {
	if _, ok := req.Get[key]; !ok {
		err := errors.New(key + " 参数异常")
		common.Error(err)
		return err
	}
	return nil
}
```