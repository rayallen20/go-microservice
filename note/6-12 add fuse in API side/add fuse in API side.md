# add fuse in API side

如前文介绍,微服务共3层:

- API网关层
- API层
- 基础服务层

上一节开发了基础服务层,这一节来开发API层和API网关层.

## PART1. 初始化项目

`docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=5A2A1531917A4D2B cap1573/cap-micro new --type=api git.imooc.com/cap1573/cartApi`

注意:此处在new后边多了一个`--type=api`

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── generate.go
├── go.mod
├── handler
│   └── cartApi.go
├── main.go
├── plugin.go
└── proto
    └── cartApi
        └── cartApi.proto

3 directories, 9 files
```

可以看到,API层和基础服务层是完全独立的

## PART2. 编写proto并生成代码

API层只需编写proto,并根据proto编写handler层即可

- step1. 编写proto

`cartApi/proto/cartApi/cartApi.proto`:

```proto
syntax = "proto3";

package go.micro.api.cartApi;


service cartApi {
	// FindAll 查询所有购物车
	rpc FindAll(Request) returns (Response) {}
}

// Pair 表示键值对的对象
message Pair {
	string key = 1;
	repeated string values = 2;
}

// Request 请求对象
message Request {
	// method 请求方法
	string method = 1;
	// path 请求路径
	string path = 2;
	// header 请求头
	map<string, Pair> header = 3;
	// get GET请求参数
	map<string, Pair> get = 4;
	// post POST请求参数
	map<string, Pair> post = 5;
	// 请求体
	string body = 6;
	// 请求的URL
	string url = 7;
}

// Response 响应对象
message Response {
	// statusCode 响应码
	int32 statusCode = 1;
	// header 响应头
	map<string, Pair> header = 2;
	// body 响应体
	string body = 3;
}
```

- step2. 生成代码

```
make proto
```

或:`sudo docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=5A2A1531917A4D2B cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/cartApi/cartApi.proto`

均可

- step3. 初始化依赖

```
go mod tidy
...
go: found github.com/micro/go-micro/v2/api in github.com/micro/go-micro/v2 v2.9.1
go: found github.com/micro/go-micro/v2/client in github.com/micro/go-micro/v2 v2.9.1
go: found github.com/micro/go-micro/v2/server in github.com/micro/go-micro/v2 v2.9.1
```

## PART3. 拉取cart模块

`go get github.com/rayallen20/cart`

模块移植到github的注意事项:

1. go.mod要写成github的地址
2. 项目中所有文件的import要改
3. 改完了之后要重新go mod tidy
4. 最后推送

## PART4. 编写handler层

和之前一样,handler层需要实现proto生成的接口.但API层的handler作用,实际上是调用基础服务.

`cartApi/handler/cartApi.go`:

```go
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	cart "github.com/rayallen20/cart/proto/cart"
	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
	"strconv"
)

type CartApi struct {
	CartService cart.CartService
}

// FindAll 通过API向外暴露为/cartApi/findAll，接收http请求
// 即：/cartApi/findAll请求会调用go.micro.api.cartApi服务的CartApi.FindAll方法
func (c *CartApi) FindAll(ctx context.Context, req *cartApi.Request, rsp *cartApi.Response) error {
	log.Info("接收到 /cartApi/findAll 访问请求")
	_, ok := req.Get["user_id"]
	if !ok {
		return errors.New("参数异常")
	}

	// 获取user_id并转化其类型
	userIdString := req.Get["user_id"].Values[0]
	fmt.Println(userIdString)
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return err
	}

	// 获取购物车中所有商品 即调用cart.CartService
	cartAll, err := c.CartService.GetAll(context.TODO(), &cart.CartFindAll{UserId: userId})

	// 将响应的结构体转化为JSON格式
	getAllBytes, err := json.Marshal(cartAll)
	if err != nil {
		return err
	}

	rsp.StatusCode = 200
	rsp.Body = string(getAllBytes)
	return nil
}
```

## PART5. 修改`main()`函数

### 5.1 添加注册中心

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/cartApi/handler"
	"github.com/rayallen20/common"

	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
)

func main() {
	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 创建服务
	svc := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		//添加 consul 注册中心
		micro.Registry(consulRegistry),
	)

	// 初始化服务
	svc.Init()

	// 注册handler
	cartApi.RegisterCartApiHandler(svc.Server(), new(handler.CartApi))

	// 运行服务
	if err = svc.Run(); err != nil {
		log.Fatal(err)
	}
}

```

API层可以不加配置中心,因为在本例中,配置中心中存储的是MySQL连接相关的参数,这些参数API层用不到.

### 5.2 添加链路追踪

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/cartApi/handler"
	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
	"github.com/rayallen20/common"
)

func main() {
	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 获取链路追踪
	t, io, err := common.NewTracer("go.micro.api.cartApi", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 创建服务
	svc := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		//添加 consul 注册中心
		micro.Registry(consulRegistry),
		//添加链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
	)

	// 初始化服务
	svc.Init()

	// 注册handler
	cartApi.RegisterCartApiHandler(svc.Server(), new(handler.CartApi))

	// 运行服务
	if err = svc.Run(); err != nil {
		log.Fatal(err)
	}
}
```
