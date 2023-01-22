# use jaeger

注:从上一章的`common`目录下将`config.go`(初始化consul配置的)和`mysql.go`(从consul中读取mysql配置的)复制过来

## PART1. 从配置中心中获取配置项并将微服务注册到注册中心

## 1.1 从配置中心中获取配置项

和上一章中对consul的调用相同,需要在`main.go`中和配置中心通信,获取项目依赖的配置项(本例中依旧是mysql相关的配置)

`product/main.go`:

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	product.RegisterProductHandler(service.Server(), new(handler.Product))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## 1.2 将微服务注册到注册中心

此处也是和上一章相同,需要将微服务注册到注册中心

`product/main.go`:

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8083"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
	)

	// Initialise service
	service.Init()

	// Register Handler
	product.RegisterProductHandler(service.Server(), new(handler.Product))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART2. 使用链路追踪

### 2.1 创建jaeger客户端实例

在`product/common`下创建`jaeger.go`:

`product/common/jaeger.go`:

```go
package common

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

// NewTracer 创建jaeger的客户端实例
func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	// TODO: 这里的配置项都是啥意思?
	conf := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: 1 * time.Second,
			LogSpans:            true,
			// jaeger接收数据的端口
			LocalAgentHostPort:  addr,
		},
	}

	tracer, closer, err := conf.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}
```

TODO: 此处的配置项意义不明,后续需要查

### 2.2 使用jaeger客户端进行链路追踪

`product/main.go`:

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 使用jaeger客户端进行链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.product", "127.0.0.1:6831")
	if err != nil {
		log.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8083"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
	)

	// Initialise service
	service.Init()

	// Register Handler
	product.RegisterProductHandler(service.Server(), new(handler.Product))

	// Run service
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

注:需实现执行`go get github.com/uber/jaeger-client-go`

## PART3. 使用配置中心中的配置项连接数据库并初始化表结构

`product/main.go`:

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 使用jaeger客户端进行链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.product", "127.0.0.1:6831")
	if err != nil {
		log.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 使用配置中心中的配置项连接数据库
	mysqlConf, err := common.GetMySQLConf(conf, "mysql")
	if err != nil {
		log.Error(err)
	}

	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@/" + mysqlConf.DataBase + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	// 禁止副表
	db.SingularTable(true)

	// 初始化 仅执行1次
	err = repository.NewProductRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8083"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
	)

	// Initialise service
	service.Init()

	// Register Handler
	product.RegisterProductHandler(service.Server(), new(handler.Product))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART4. 创建service层的实例

此处的service层指的是`product/domain/service.go`中的实例

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/domain/repository"
	service2 "git.imooc.com/cap1573/product/domain/service"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 使用jaeger客户端进行链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.product", "127.0.0.1:6831")
	if err != nil {
		log.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 使用配置中心中的配置项连接数据库
	mysqlConf, err := common.GetMySQLConf(conf, "mysql")
	if err != nil {
		log.Error(err)
	}

	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@/" + mysqlConf.DataBase + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	// 禁止副表
	db.SingularTable(true)

	// 初始化 仅执行1次
	err = repository.NewProductRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// 创建service层的实例
	productDataService := service2.NewProductDataService(repository.NewProductRepository(db))

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8083"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
	)

	// Initialise service
	service.Init()

	// 注册Handler
	product.RegisterProductHandler(service.Server(), &handler.Product{ProductDataService: productDataService})

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## PART5. 设置服务

微服务名、版本、微服务监听地址和端口、将微服务添加到注册中心这几步都已经完成了.此处需要绑定链路追踪

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/domain/repository"
	service2 "git.imooc.com/cap1573/product/domain/service"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 使用jaeger客户端进行链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.product", "127.0.0.1:6831")
	if err != nil {
		log.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 使用配置中心中的配置项连接数据库
	mysqlConf, err := common.GetMySQLConf(conf, "mysql")
	if err != nil {
		log.Error(err)
	}

	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@/" + mysqlConf.DataBase + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	// 禁止副表
	db.SingularTable(true)

	// 初始化 仅执行1次
	err = repository.NewProductRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// 创建service层的实例
	productDataService := service2.NewProductDataService(repository.NewProductRepository(db))

	// 设置服务
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8083"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		// TODO: 此处的micro.WrapClient是干啥的? 和micro.WrapHandler有啥区别?
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// 初始化服务
	service.Init()

	// 注册Handler
	err = product.RegisterProductHandler(service.Server(), &handler.Product{ProductDataService: productDataService})
	if err != nil {
		log.Error(err)
	}

	// 运行服务
	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

TODO:此处的`micro.WrapClient`是干啥的? 和`micro.WrapHandler`有啥区别?

注:需事先执行`go get github.com/micro/go-plugins/wrapper/trace/opentracing/v2`

## PART6. 完整的`main()`函数

```go
package main

import (
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/domain/repository"
	service2 "git.imooc.com/cap1573/product/domain/service"
	"git.imooc.com/cap1573/product/handler"
	product "git.imooc.com/cap1573/product/proto/product"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 从配置中心获取配置项
	consulHost := "127.0.0.1"
	consulPort := 8500
	keyPrefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(consulHost, int64(consulPort), keyPrefix)
	if err != nil {
		log.Error(err)
	}

	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 使用jaeger客户端进行链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.product", "127.0.0.1:6831")
	if err != nil {
		log.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 使用配置中心中的配置项连接数据库
	mysqlConf, err := common.GetMySQLConf(conf, "mysql")
	if err != nil {
		log.Error(err)
	}

	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@/" + mysqlConf.DataBase + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	// 禁止副表
	db.SingularTable(true)

	// 初始化 仅执行1次
	err = repository.NewProductRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// 创建service层的实例
	productDataService := service2.NewProductDataService(repository.NewProductRepository(db))

	// 设置服务
	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8083"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		// TODO: 此处的micro.WrapClient是干啥的?
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// 初始化服务
	service.Init()

	// 注册Handler
	err = product.RegisterProductHandler(service.Server(), &handler.Product{ProductDataService: productDataService})
	if err != nil {
		log.Error(err)
	}

	// 运行服务
	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

推送代码:

```
git init
git add .
git commit -m "product微服务开发完成"
git push https://我的慕课git账户:我的慕课git密码@git.imooc.com/rayallen20c/product.git --all
```

## PART7. 运行product微服务时的consul

![category和product同时运行时的consul](./img/category和product同时运行时的consul.png)

## PART8. 编写测试代码

注:正式使用时,需要先`go get git.imooc.com/rayallen20c/product`

此时微服务product已经配置了链路追踪,还需要有一个同样配置了链路追踪的微服务作为它的客户端去请求它,才能在jaeger中看到链路追踪的span.

在`product`目录下创建`client.go`.

注:需要将`product/proto/product/product.proto`中的message`ProductSeo.seo_keyword`更改为`ProductSeo.seo_keywords`,然后重新生成go文件.因为之前写的一个字段是错误的

`product/client.go`:

```go
package main

import (
	"context"
	"fmt"
	"git.imooc.com/cap1573/product/common"
	go_micro_service_product "git.imooc.com/cap1573/product/proto/product"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 使用jaeger客户端进行链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.product.client", "127.0.0.1:6831")
	if err != nil {
		log.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 设置服务
	service := micro.NewService(
		micro.Name("go.micro.service.product.client"),
		micro.Version("latest"),
		// 设置商品微服务监听的地址和端口
		// Tips: 若要配合docker和K8S 此处需更改为监听0.0.0.0
		micro.Address("127.0.0.1:8084"),
		// 将微服务注册到注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		// 此处作为客户端的微服务就要写WrapClient
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// TODO:若一个微服务既要请求其他微服务 又要作为其他微服务的服务端 是否这俩都配上?
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// 创建product微服务
	productService := go_micro_service_product.NewProductService("go.micro.service.product", service.Client())

	productAdd := &go_micro_service_product.ProductInfo{
		ProductName:        "testName1",
		ProductSku:         "testSku1",
		ProductPrice:       1.5,
		ProductDescription: "测试商品描述",
		ProductCategoryId:  1,
		ProductImage: []*go_micro_service_product.ProductImage{
			{
				ImageName: "测试图片名称1",
				ImageCode: "testImageCode1",
				ImageUrl:  "testImageUrl1",
			},
			{
				ImageName: "测试图片名称2",
				ImageCode: "testImageCode2",
				ImageUrl:  "testImageUrl2",
			},
		},
		ProductSize: []*go_micro_service_product.ProductSize{
			{
				SizeName: "testSizeName",
				SizeCode: "testSizeCode",
			},
		},
		ProductSeo: &go_micro_service_product.ProductSeo{
			SeoTitle:       "testSeoTitle",
			SeoKeywords:    "testSeoKeywords",
			SeoDescription: "testSeoDescription",
			SeoCode:        "testSeoCode",
		},
	}

	// 调用微服务中的AddProduct API
	response, err := productService.AddProduct(context.TODO(), productAdd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}
```

运行后查看jaeger:

![运行后查看链路追踪1.png](./img/运行后查看链路追踪1.png)

![运行后查看链路追踪2.png](./img/运行后查看链路追踪2.png)

注意:**本例中调用的微服务是在本地中已经存在了,若本地不存在要调用的微服务,需要将这个被调用的微服务`go get`下来**


























