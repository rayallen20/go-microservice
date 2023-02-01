# add current-limiting server side

## PART1. 添加配置中心

`cart/main.go`

```go
package main

import (
	"git.imooc.com/cap1573/cart/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/rayallen20/common"

	cart "git.imooc.com/cap1573/cart/proto/cart"
)

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}
	
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	cart.RegisterCartHandler(service.Server(), new(handler.Cart))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART2. 添加注册中心

```go
package main

import (
	"git.imooc.com/cap1573/cart/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/rayallen20/common"

	cart "git.imooc.com/cap1573/cart/proto/cart"
)

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}

	// 将自身注册到注册中心
	register := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	cart.RegisterCartHandler(service.Server(), new(handler.Cart))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART3. 添加链路追踪

```go
package main

import (
	"git.imooc.com/cap1573/cart/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"

	cart "git.imooc.com/cap1573/cart/proto/cart"
)

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}

	// 将自身注册到注册中心
	register := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	cart.RegisterCartHandler(service.Server(), new(handler.Cart))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART4. 添加数据库连接

```go
package main

import (
	"git.imooc.com/cap1573/cart/domain/repository"
	"git.imooc.com/cap1573/cart/handler"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"

	cart "git.imooc.com/cap1573/cart/proto/cart"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}

	// 将自身注册到注册中心
	register := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 获取数据库连接所需配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		log.Error(err)
	}
	// 创建数据库连接
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	// 禁止复表
	db.SingularTable(true)
	// 初始化数据库
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	cart.RegisterCartHandler(service.Server(), new(handler.Cart))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

注:此处不要忘记import `_ "github.com/jinzhu/gorm/dialects/mysql"`

## PART5. 指定服务监听端口、指定注册中心、启用链路追踪

```go
package main

import (
	"git.imooc.com/cap1573/cart/domain/repository"
	"git.imooc.com/cap1573/cart/handler"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"

	cart "git.imooc.com/cap1573/cart/proto/cart"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}

	// 将自身注册到注册中心
	register := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 获取数据库连接所需配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		log.Error(err)
	}
	// 创建数据库连接
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	// 禁止复表
	db.SingularTable(true)
	// 初始化数据库
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		// 指定服务端口
		micro.Address("0.0.0.0:8087"),
		// 指定注册中心
		micro.Registry(register),
		// 启用链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// Initialise service
	service.Init()

	// Register Handler
	cart.RegisterCartHandler(service.Server(), new(handler.Cart))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART6. 添加限流

再次注意:限流在服务端做,熔断在客户端做

执行如下命令:`go get github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2`

```go
package main

import (
	"git.imooc.com/cap1573/cart/domain/repository"
	"git.imooc.com/cap1573/cart/handler"
	cart "git.imooc.com/cap1573/cart/proto/cart"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
)

// QPS 每秒可处理的最大请求数量 以该值作为限流的阈值 超过该值部分的流量将被丢弃
var QPS = 100

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}

	// 实例化注册中心
	register := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 获取数据库连接所需配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		log.Error(err)
	}
	// 创建数据库连接
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	// 禁止复表
	db.SingularTable(true)
	// 初始化数据库
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		// 指定服务端口
		micro.Address("0.0.0.0:8087"),
		// 指定注册中心
		micro.Registry(register),
		// 启用链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	service.Init()

	// Register Handler
	cart.RegisterCartHandler(service.Server(), new(handler.Cart))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART7. 服务初始化

```go
package main

import (
	"git.imooc.com/cap1573/cart/domain/repository"
	service2 "git.imooc.com/cap1573/cart/domain/service"
	"git.imooc.com/cap1573/cart/handler"
	cart "git.imooc.com/cap1573/cart/proto/cart"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
)

// QPS 每秒可处理的最大请求数量 以该值作为限流的阈值 超过该值部分的流量将被丢弃
var QPS = 100

func main() {
	// 从配置中心读取配置
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/my-config")
	if err != nil {
		log.Error(err)
	}

	// 实例化注册中心
	register := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 获取数据库连接所需配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		log.Error(err)
	}
	// 创建数据库连接
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	// 禁止复表
	db.SingularTable(true)
	// 初始化数据库
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		// 指定服务端口
		micro.Address("0.0.0.0:8087"),
		// 指定注册中心
		micro.Registry(register),
		// 启用链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	service.Init()

	// 初始化service层
	cartDataService := service2.NewCartDataService(repository.NewCartRepository(db))

	// Register Handler
	cart.RegisterCartHandler(service.Server(), &handler.Cart{CartDataService: cartDataService})

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART8. 推送代码

```
git init
git add .
git commit -m "cart微服务开发完成"
git push https://慕课网账户:慕课网密码@git.imooc.com/rayallen20c/cart.git --all
```
