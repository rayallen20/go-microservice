# develop main--first half

## PART1. 注册中心

`paymentApi/main.go`:

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/rayallen20/common"
	"github.com/rayallen20/paymentApi/handler"

	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
)

func main() {
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.paymentApi"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	paymentApi.RegisterPaymentApiHandler(service.Server(), new(handler.PaymentApi))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART2. 链路追踪

`paymentApi/main.go`:

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/paymentApi/handler"

	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
)

func main() {
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.api.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.paymentApi"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	paymentApi.RegisterPaymentApiHandler(service.Server(), new(handler.PaymentApi))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART3. 熔断

`paymentApi/main.go`:

```go
package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/paymentApi/handler"
	"net"
	"net/http"

	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
)

func main() {
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.api.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 熔断启动监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9192"), hystrixStreamHandler)
		if err != nil {
			common.Error(err)
			log.Error(err)
		}
	}()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.paymentApi"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	paymentApi.RegisterPaymentApiHandler(service.Server(), new(handler.PaymentApi))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART4. 监控

`paymentApi/main.go`:

```go
package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/paymentApi/handler"
	"net"
	"net/http"

	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
)

func main() {
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.api.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 熔断启动监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9192"), hystrixStreamHandler)
		if err != nil {
			common.Error(err)
			log.Error(err)
		}
	}()
	
	// 监控
	common.BootPrometheus(9292)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.paymentApi"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	paymentApi.RegisterPaymentApiHandler(service.Server(), new(handler.PaymentApi))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART5. 创建并初始化服务

`paymentApi/main.go`:

```go
package main

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/paymentApi/handler"
	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
	"net"
	"net/http"
)

func main() {
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.api.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 熔断启动监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9192"), hystrixStreamHandler)
		if err != nil {
			common.Error(err)
			log.Error(err)
		}
	}()

	// 监控
	common.BootPrometheus(9292)

	// 创建服务
	paymentApiService := micro.NewService(
		micro.Name("go.micro.api.paymentApi"),
		micro.Version("latest"),
		// 注册中心
		micro.Registry(consulRegistry),
		// 链路追踪
		// 作为服务端访问的链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 作为客户端访问的链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		// 负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// 初始化服务
	paymentApiService.Init()

	// Register Handler
	paymentApi.RegisterPaymentApiHandler(paymentApiService.Server(), new(handler.PaymentApi))

	// Run service
	if err = paymentApiService.Run(); err != nil {
		log.Fatal(err)
	}
}

// clientWrapper 作为熔断器的客户端结构体使用
type clientWrapper struct {
	client.Client
}

// Call 方法签名要和client.Client.Call相同
// 以便实现接口go-micro/v2/client.Client
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// 以请求的服务名称 和 请求的IP地址作为熔断器中计数器的名称
	// 这个名称用于让计数器在记录状态(调用次数,失败次数,被拒绝次数等信息)时区分具体哪一个服务
	// 加IP是为了在有多个同样的微服务部署在不同容器中的场景下做区分
	name := req.Service() + "." + req.Endpoint()

	// 正常状态下执行的函数
	runFunc := func() error {
		common.Info(name)
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// 熔断状态下执行的函数
	fallbackFunc := func(err error) error {
		common.Error(err)
		return err
	}
	return hystrix.Do(name, runFunc, fallbackFunc)
}

// NewClientHystrixWrapper 创建熔断器客户端
func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
```

## PART6. 注册handler并启动服务

`paymentApi/main.go`:

```go
package main

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	go_micro_service_payment "github.com/rayallen20/payment/proto/payment"
	"github.com/rayallen20/paymentApi/handler"
	paymentApi "github.com/rayallen20/paymentApi/proto/paymentApi"
	"net"
	"net/http"
)

func main() {
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.api.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	// 熔断启动监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9192"), hystrixStreamHandler)
		if err != nil {
			common.Error(err)
			log.Error(err)
		}
	}()

	// 监控
	common.BootPrometheus(9292)

	// 创建服务
	paymentApiService := micro.NewService(
		micro.Name("go.micro.api.paymentApi"),
		micro.Version("latest"),
		// 注册中心
		micro.Registry(consulRegistry),
		// 链路追踪
		// 作为服务端访问的链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 作为客户端访问的链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		// 负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// 初始化服务
	paymentApiService.Init()

	// 注册Handler
	paymentService := go_micro_service_payment.NewPaymentService("go.micro.service.payment", paymentApiService.Client())
	err = paymentApi.RegisterPaymentApiHandler(paymentApiService.Server(), &handler.PaymentApi{PaymentService: paymentService})
	if err != nil {
		common.Error(err)
	}

	// 启动服务
	if err = paymentApiService.Run(); err != nil {
		log.Fatal(err)
	}
}

// clientWrapper 作为熔断器的客户端结构体使用
type clientWrapper struct {
	client.Client
}

// Call 方法签名要和client.Client.Call相同
// 以便实现接口go-micro/v2/client.Client
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// 以请求的服务名称 和 请求的IP地址作为熔断器中计数器的名称
	// 这个名称用于让计数器在记录状态(调用次数,失败次数,被拒绝次数等信息)时区分具体哪一个服务
	// 加IP是为了在有多个同样的微服务部署在不同容器中的场景下做区分
	name := req.Service() + "." + req.Endpoint()

	// 正常状态下执行的函数
	runFunc := func() error {
		common.Info(name)
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// 熔断状态下执行的函数
	fallbackFunc := func(err error) error {
		common.Error(err)
		return err
	}
	return hystrix.Do(name, runFunc, fallbackFunc)
}

// NewClientHystrixWrapper 创建熔断器客户端
func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
```