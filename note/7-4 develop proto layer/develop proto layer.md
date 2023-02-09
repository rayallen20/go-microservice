# develop proto layer

## PART1. 微服务订单领域开发

- 订单需求分析 & 项目目录创建
- 订单代码开发
- 订单系统接入Pormetheus监控

## PART2. 创建项目

```
docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=5A2A1531917A4D2B cap1573/cap-micro new github.com/rayallen20/order
```

整理后项目目录结构如下:

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
│   │   └── order.go
│   ├── repository
│   │   └── order_repository.go
│   └── service
│       └── order_data_service.go
├── generate.go
├── go.mod
├── handler
│   └── order.go
├── main.go
├── plugin.go
└── proto
    └── order
        └── order.proto

8 directories, 13 files
```

## PART3. 开发proto

`order/proto/order/order.proto`:

```proto
syntax = "proto3";

package go.micro.service.order;

service Order {
	// GetOrderByID 根据订单ID查找订单信息
	rpc GetOrderByID(OrderID) returns (OrderInfo) {}
	// GetAllOrder 查询所有订单
	rpc GetAllOrder(AllOrderRequest) returns (AllOrder) {}
	// CreateOrder 创建订单
	rpc CreateOrder(OrderInfo) returns (OrderID) {}
	// DeleteOrderByID 根据ID删除订单信息
	rpc DeleteOrderByID(OrderID) returns (Response) {}
	// UpdateOrderPayStatus 更新订单支付状态
	rpc UpdateOrderPayStatus(PayStatus) returns (Response) {}
	// UpdateOrderShipStatus 更新订单发货状态
	rpc UpdateOrderShipStatus(ShipStatus) returns (Response) {}
	// UpdateOrder 更新订单信息
	rpc UpdateOrder(OrderInfo) returns (Response) {}
}

// AllOrderRequest 查询所有订单的请求结构体
message AllOrderRequest {

}

// AllOrder 查询所有订单的响应
message AllOrder {
	// order_info 订单信息集合
	repeated OrderInfo order_info = 1;
}

// OrderID 订单ID信息
message OrderID {
	int64 order_id = 1;
}

// OrderInfo 订单信息
message OrderInfo {
	// id 订单id
	int64 id = 1;
	// pay_status 订单支付状态
	int32 pay_status = 2;
	// ship_status 订单发货状态
	int32 ship_status = 3;
	// price 订单总价
	double price = 4;
	// order_detail 订单详情信息集合
	repeated OrderDetail order_detail = 5;
}

// OrderDetail 订单详情信息
message OrderDetail {
	// id 订单详情id
	int64 id = 1;
	// product_id 商品信息id
	int64 product_id = 2;
	// product_num 商品数量
	int64 product_num = 3;
	// product_size_id 商品尺码
	int64 product_size_id = 4;
	// product_price 商品价格
	int64 product_price = 5;
	// order_id 订单详情所属的订单id
	int64 order_id = 6;
}

// Response 响应结构体
message Response {
	string msg = 1;
}

// PayStatus 订单支付状态结构体
message PayStatus {
	// order_id 订单ID
	int64 order_id = 1;
	// pay_status 订单支付状态
	int32 pay_status = 2;
}

// ShipStatus 订单发货状态结构体
message ShipStatus {
	// order_id 订单ID
	int64 order_id = 1;
	// ship_status 订单发货状态
	int32 ship_status = 2;
}
```

## PART4. 生成代码

根目录下执行`make proto`即可

执行后目录结构如下:

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
│   │   └── order.go
│   ├── repository
│   │   └── order_repository.go
│   └── service
│       └── order_data_service.go
├── generate.go
├── go.mod
├── handler
│   └── order.go
├── main.go
├── plugin.go
└── proto
    └── order
        ├── order.pb.go
        ├── order.pb.micro.go
        └── order.proto

8 directories, 15 files
```