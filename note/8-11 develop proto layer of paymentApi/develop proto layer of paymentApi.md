# develop proto layer of paymentApi

## PART1. 初始化项目

项目根目录下执行:

`docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=5A2A1531917A4D2B cap1573/cap-micro new --type=api github.com/rayallen20/paymentApi`

整理后的目录结构如下:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── generate.go
├── go.mod
├── handler
│   └── paymentApi.go
├── main.go
├── plugin.go
└── proto
    └── paymentApi
        └── paymentApi.proto

3 directories, 9 files
```

## PART2. 编写proto

`paymentApi/proto/paymentApi/paymentApi.proto`:

```proto
syntax = "proto3";

package go.micro.api.paymentApi;


service PaymentApi {
	// PayPalRefund 直接退款功能
	rpc PayPalRefund(Request) returns (Response) {}
}

message Pair {
	string key = 1;
	repeated string values = 2;
}


message Request {
	string method = 1;
	string path = 2;
	map<string, Pair> header = 3;
	map<string, Pair> get = 4;
	map<string, Pair> post = 5;
	string body = 6;
	string url = 7;
}


message Response {
	int32 statusCode = 1;
	map<string, Pair> header = 2;
	string body = 3;
}
```

项目根目录下执行:`make proto`

执行完毕后目录结构如下:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── generate.go
├── go.mod
├── handler
│   └── paymentApi.go
├── main.go
├── plugin.go
└── proto
    └── paymentApi
        ├── paymentApi.pb.go
        ├── paymentApi.pb.micro.go
        └── paymentApi.proto

3 directories, 11 files
```