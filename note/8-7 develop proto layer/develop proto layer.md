# develop proto layer

## PART1. 初始化工程

`docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=5A2A1531917A4D2B cap1573/cap-micro new github.com/rayallen20/payment`

整理后目录结构如下:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── common
│   └── README.md
├── domain
│   ├── model
│   │   └── payment.go
│   ├── repository
│   │   └── payment_repository.go
│   └── service
│       └── payment_data_service.go
├── generate.go
├── go.mod
├── handler
│   └── payment.go
├── main.go
├── plugin.go
└── proto
    └── payment
        └── payment.proto

8 directories, 13 files
```

## PART2. 编写`payment.proto`

`payment/proto/payment/payment.proto`:

```proto
syntax = "proto3";

package go.micro.service.payment;

service Payment {
	// AddPayment 创建支付单
	rpc AddPayment(PaymentInfo) returns (PaymentID) {}
	// UpdatePayment 更新支付单信息
	rpc UpdatePayment(PaymentInfo) returns (Response) {}
	// DeletePayment 删除支付单信息
	rpc DeletePayment(PaymentID) returns (Response) {}
	// FindPaymentByID 根据ID查找支付单信息
	rpc FindPaymentByID(PaymentID) returns (Response) {}
	// FindAllPayment 查找所有支付单信息
	rpc FindAllPayment(All) returns (PaymentAll) {}
}

// PaymentInfo 支付单信息
message PaymentInfo {
	// id 支付单id
	int64 id = 1;
	// payment_name 支付单名称
	string payment_name = 2;
	// paypal应用的secret id
	string payment_sid = 3;
	// payment_status 支付环境状态
	// false: sandbox状态
	// true: live状态
	bool payment_status = 4;
	// payment_image 支付单图片
	string payment_image = 5;
}

// PaymentID 支付单id信息
message PaymentID {
	// payment_id 支付单id
	int64 payment_id = 1;
}

// Response 响应信息
message Response {
	// msg 响应信息
	string msg = 1;
}

// All 查找所有支付单信息参数
message All {

}

// PaymentAll 支付单信息集合
message PaymentAll {
	// payment_info 支付单信息集合
	repeated PaymentInfo payment_info = 1;
}
```

在项目根目录下:`make proto`

```
make proto
恭喜，恭喜命令执行成功！
```

执行完毕后的目录结构:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── common
│   └── README.md
├── domain
│   ├── model
│   │   └── payment.go
│   ├── repository
│   │   └── payment_repository.go
│   └── service
│       └── payment_data_service.go
├── generate.go
├── go.mod
├── handler
│   └── payment.go
├── main.go
├── plugin.go
└── proto
    └── payment
        ├── payment.pb.go
        ├── payment.pb.micro.go
        └── payment.proto

8 directories, 15 files
```