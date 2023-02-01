# add load balance in API side

## PART1. 添加熔断

再次强调:熔断是加在客户端的,也就是API层;限流是加在服务端的,也就是基础服务层.

拉取hystrix包:`go get github.com/afex/hystrix-go/hystrix`

### 1.1 创建、启动熔断器并监听

```go
package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/cartApi/handler"
	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
	"github.com/rayallen20/common"
	"net"
	"net/http"
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

	// 添加熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	// 启动熔断器
	hystrixStreamHandler.Start()
	// 熔断器监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9096"), hystrixStreamHandler)
		if err != nil {
			log.Error(err)
		}
	}()

	// 创建服务
	svc := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		// 指定微服务地址和端口
		micro.Address("0.0.0.0:8086"),
		// 添加 consul 注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		// NewClientWrapper:作为客户端访问时 使用该函数添加链路追踪
		// NewHandlerWrapper:作为服务端被访问时 使用该函数添加链路追踪
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

### 1.2 创建熔断器的客户端

```go
package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/cartApi/handler"
	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
	"github.com/rayallen20/common"
	"net"
	"net/http"
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

	// 创建熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	// 启动熔断器
	hystrixStreamHandler.Start()
	// 熔断器监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9096"), hystrixStreamHandler)
		if err != nil {
			log.Error(err)
		}
	}()

	// 创建服务
	svc := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		// 指定微服务地址和端口
		micro.Address("0.0.0.0:8086"),
		// 添加 consul 注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		// NewClientWrapper:作为客户端访问时 使用该函数添加链路追踪
		// NewHandlerWrapper:作为服务端被访问时 使用该函数添加链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
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
		fmt.Println(name)
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// 熔断状态下执行的函数
	fallbackFunc := func(err error) error {
		fmt.Println(err)
		return err
	}
	return hystrix.Do(name, runFunc, fallbackFunc)
}

// NewClientHystrixWrapper 创建熔断器客户端
func NewClientHystrixWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
```

## PART2. 添加负载均衡

负载均衡也是在客户端添加的.

拉取roundrobin包:`go get github.com/micro/go-plugins/wrapper/select/roundrobin/v2`

```go
package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/cartApi/handler"
	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
	"github.com/rayallen20/common"
	"net"
	"net/http"
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

	// 创建熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	// 启动熔断器
	hystrixStreamHandler.Start()
	// 熔断器监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9096"), hystrixStreamHandler)
		if err != nil {
			log.Error(err)
		}
	}()

	// 创建服务
	svc := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		// 指定微服务地址和端口
		micro.Address("0.0.0.0:8086"),
		// 添加 consul 注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		// NewClientWrapper:作为客户端访问时 使用该函数添加链路追踪
		// NewHandlerWrapper:作为服务端被访问时 使用该函数添加链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		// 添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
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
		fmt.Println(name)
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// 熔断状态下执行的函数
	fallbackFunc := func(err error) error {
		fmt.Println(err)
		return err
	}
	return hystrix.Do(name, runFunc, fallbackFunc)
}

// NewClientHystrixWrapper 创建熔断器客户端
func NewClientHystrixWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
```

## PART3. 创建要访问的服务并初始化

```go
package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	cart "github.com/rayallen20/cart/proto/cart"
	"github.com/rayallen20/cartApi/handler"
	cartApi "github.com/rayallen20/cartApi/proto/cartApi"
	"github.com/rayallen20/common"
	"net"
	"net/http"
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

	// 创建熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	// 启动熔断器
	hystrixStreamHandler.Start()
	// 熔断器监听
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9096"), hystrixStreamHandler)
		if err != nil {
			log.Error(err)
		}
	}()

	// 创建服务
	svc := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		// 指定微服务地址和端口
		micro.Address("0.0.0.0:8086"),
		// 添加 consul 注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		// NewClientWrapper:作为客户端访问时 使用该函数添加链路追踪
		// NewHandlerWrapper:作为服务端被访问时 使用该函数添加链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		// 添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// 初始化服务
	svc.Init()

	// 创建要访问的微服务
	// Tips: 这个NewCartService是基础服务cart的proto生成的
	cartService := cart.NewCartService("go.micro.service.cart", svc.Client())

	// 注册handler
	err = cartApi.RegisterCartApiHandler(svc.Server(), &handler.CartApi{CartService: cartService})
	if err != nil {
		log.Error(err)
	}

	// 运行服务
	if err = svc.Run(); err != nil {
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
		fmt.Println(name)
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// 熔断状态下执行的函数
	fallbackFunc := func(err error) error {
		fmt.Println(err)
		return err
	}
	return hystrix.Do(name, runFunc, fallbackFunc)
}

// NewClientHystrixWrapper 创建熔断器客户端
func NewClientHystrixWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
```