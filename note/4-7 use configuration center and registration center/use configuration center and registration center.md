# use configuration center and registration center


## PART1. 实现配置中心

### 1.1 实现配置中心

由于配置中心是每一个微服务都需要用到的,所以将初始化配置中心的代码写在`user/git.imooc.com/cap1573/category/common`中.

`user/git.imooc.com/cap1573/category/common/config.go`:

```go
package common

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-plugins/config/source/consul/v2"
	"strconv"
)

// GetConsulConfig 初始化Consul的配置
func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
	address := host + ":" + strconv.FormatInt(port, 10)
	consulSource := consul.NewSource(
		consul.WithAddress(address),
		// 不配置前缀的话 默认前缀为 /micro/config
		consul.WithPrefix(prefix),
		// 设置是否移除前缀 设置为true表示可以在没有前缀的前提下获取对应配置
		consul.StripPrefix(true),
	)

	// 初始化配置
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	err = conf.Load(consulSource)
	if err != nil {
		return nil, err
	}

	return conf, err
}
```

注:需要执行:`go get github.com/micro/go-plugins/config/source/consul/v2`,加载go-micro的插件模块

### 1.2 使用配置中心

在`main.go`中调用该函数即可:

```go
package main

import (
	"git.imooc.com/cap1573/category/common"
	"git.imooc.com/cap1573/category/handler"
	category "git.imooc.com/cap1573/category/proto/category"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// 获取配置中心
	host := "127.0.0.1"
	port := 8500
	prefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(host, int64(port), prefix)
	if err != nil {
		log.Error(err)
	}
	
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	category.RegisterCategoryHandler(service.Server(), new(handler.Category))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

注:此处除获取配置中心部分的代码外,剩下的都是脚手架自动生成的.

## PART2. 实现注册中心

直接在`main.go`中创建注册中心即可

```go
package main

import (
	"git.imooc.com/cap1573/category/common"
	"git.imooc.com/cap1573/category/handler"
	category "git.imooc.com/cap1573/category/proto/category"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func main() {
	// 获取配置中心
	host := "127.0.0.1"
	port := 8500
	prefix := "/micro/my-config"
	conf, err := common.GetConsulConfig(host, int64(port), prefix)
	if err != nil {
		log.Error(err)
	}
	
	// 获取注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string {
			"127.0.0.1:8500",
		}
	})
	
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	category.RegisterCategoryHandler(service.Server(), new(handler.Category))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

- 注:此处配置中心和注册中心只是使用了同一个中间件(consul),但实际上它们的作用并不相同(从编码中也可以看出,类型其实也不同).但具体有啥区别,我现在也不知道
- 注:此时代码时跑不起来的,因为变量`conf`和`consulRegistry`还未被使用
- 注:需使用`go get github.com/micro/go-plugins/registry/consul/v2`,加载注册中心的consul插件.注意这个插件和1.1章节中的配置中心插件不同
- 注:需使用`go get github.com/micro/go-micro/v2/registry`,加载注册中心插件

