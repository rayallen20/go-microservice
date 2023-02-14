# develop main

## PART1. 配置中心与注册中心

### 1.1 配置中心

`payment/main.go`:

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/handler"

	payment "github.com/rayallen20/payment/proto/payment"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), new(handler.Payment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### 1.2 注册中心

`payment/main.go`:

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), new(handler.Payment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART2. 链路追踪

`payment/main.go`:

```go
package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), new(handler.Payment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART3. 初始化数据库

`payment/main.go`:

```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/repository"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 从配置中心读取MySQL的配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		common.Error(err)
	}

	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Error(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 以下代码仅执行1次
	// 创建表
	tableInit := repository.NewPaymentRepository(db)
	err = tableInit.InitTable()
	if err != nil {
		common.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), new(handler.Payment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART4. 添加监控

`payment/main.go`:

```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/repository"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 从配置中心读取MySQL的配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		common.Error(err)
	}

	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Error(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 以下代码仅执行1次
	// 创建表
	tableInit := repository.NewPaymentRepository(db)
	err = tableInit.InitTable()
	if err != nil {
		common.Error(err)
	}

	// 添加监控
	common.BootPrometheus(9089)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), new(handler.Payment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART4. 添加监控

`payment/main.go`:

```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/repository"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 从配置中心读取MySQL的配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		common.Error(err)
	}

	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Error(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 以下代码仅执行1次
	// 创建表
	tableInit := repository.NewPaymentRepository(db)
	err = tableInit.InitTable()
	if err != nil {
		common.Error(err)
	}

	// 监控
	common.BootPrometheus(9089)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	payment.RegisterPaymentHandler(service.Server(), new(handler.Payment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART5. 创建并初始化服务

`payment/main.go`:

```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/repository"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

var (
	QPS = 1000
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 从配置中心读取MySQL的配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		common.Error(err)
	}

	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Error(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 以下代码仅执行1次
	// 创建表
	tableInit := repository.NewPaymentRepository(db)
	err = tableInit.InitTable()
	if err != nil {
		common.Error(err)
	}

	// 监控
	common.BootPrometheus(9089)

	// 创建服务
	paymentService := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
		// 指定地址与端口
		micro.Address("0.0.0.0:8089"),
		// 添加注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 加载限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		// 加载监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 初始化服务
	paymentService.Init()

	// 注册Handler
	payment.RegisterPaymentHandler(paymentService.Server(), new(handler.Payment))

	// Run service
	if err = paymentService.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART6. 注册handler

`payment/main.go`:

```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/repository"
	"github.com/rayallen20/payment/domain/service"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

var (
	QPS = 1000
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 从配置中心读取MySQL的配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		common.Error(err)
	}

	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Error(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 以下代码仅执行1次
	// 创建表
	tableInit := repository.NewPaymentRepository(db)
	err = tableInit.InitTable()
	if err != nil {
		common.Error(err)
	}

	// 监控
	common.BootPrometheus(9089)

	// 创建服务
	paymentService := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
		// 指定地址与端口
		micro.Address("0.0.0.0:8089"),
		// 添加注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 加载限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		// 加载监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 初始化服务
	paymentService.Init()

	// 注册Handler
	paymentDataService := service.NewPaymentDataService(repository.NewPaymentRepository(db))
	err = payment.RegisterPaymentHandler(paymentService.Server(), &handler.Payment{PaymentDataService: paymentDataService})
	if err != nil {
		common.Error(err)
	}

	// 运行服务
	if err = paymentService.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART7. 完整代码

完整的`main.go`如下:

```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rayallen20/common"
	"github.com/rayallen20/payment/domain/repository"
	"github.com/rayallen20/payment/domain/service"
	"github.com/rayallen20/payment/handler"
	payment "github.com/rayallen20/payment/proto/payment"
)

var (
	QPS = 1000
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/my-config")
	if err != nil {
		common.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})

	// 链路追踪
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 从配置中心读取MySQL的配置
	mysqlConf, err := common.GetMySQLConf(consulConfig, "mysql")
	if err != nil {
		common.Error(err)
	}

	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlConf.User+":"+mysqlConf.Password+"@/"+mysqlConf.DataBase+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Error(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 以下代码仅执行1次
	// 创建表
	tableInit := repository.NewPaymentRepository(db)
	err = tableInit.InitTable()
	if err != nil {
		common.Error(err)
	}

	// 监控
	common.BootPrometheus(9089)

	// 创建服务
	paymentService := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
		// 指定地址与端口
		micro.Address("0.0.0.0:8089"),
		// 添加注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 加载限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		// 加载监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 初始化服务
	paymentService.Init()

	// 注册Handler
	paymentDataService := service.NewPaymentDataService(repository.NewPaymentRepository(db))
	err = payment.RegisterPaymentHandler(paymentService.Server(), &handler.Payment{PaymentDataService: paymentDataService})
	if err != nil {
		common.Error(err)
	}

	// 运行服务
	if err = paymentService.Run(); err != nil {
		log.Fatal(err)
	}
}
```

注:此时运行,若有错误,则会在项目根目录下创建日志文件`micro.log`.