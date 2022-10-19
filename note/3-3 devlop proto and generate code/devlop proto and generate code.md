# devlop proto and generate code

## PART1. 用户功能梳理及开发

### 1.1 用户领域需求分析

- 用户信息管理
- 用户登录,注册
- 用户鉴权

### 1.2 修改模块名称

修改前的`go.mod`:

```
module user

go 1.18

require (
	github.com/golang/protobuf latest
	github.com/micro/micro/v3 latest
	google.golang.org/protobuf latest
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.27.1
```

此时模块名为`user`,给其他模块使用时,名称不友好.修改为git仓库的名称.即`git.imooc.com/rayallen20c/user`

删除该文件,然后在命令行中执行:

```
go mod init git.imooc.com/rayallen20c/user
go: creating new go.mod: module git.imooc.com/rayallen20c/user
go: to add module requirements and sums:
        go mod tidy
```

重新生成的`go.mod`内容如下:

```
module git.imooc.com/rayallen20c/user

go 1.18
```

### 1.3 下载依赖

```
go mod tidy
```

### 1.4 修改模块名

以`main.go`为例,修改前的代码为:

```go
package main

import (
	"user/handler"
	pb "user/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("user"),
	)

	// Register handler
	pb.RegisterUserHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

需要将import的`"user/handler"`,修改为`"git.imooc.com/rayallen20c/user/handler"`

## PART2. 开发用户领域功能

### 2.1 删除go-micro生成的proto文件

删除工程目录下的`/proto/user.proto`即可

### 2.2 创建领域的proto文件

- step1. 在`proto`目录下新建目录`user`

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── domain
│   ├── model
│   ├── repository
│   └── service
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── micro.mu
└── proto
    └── user

7 directories, 9 files
```

- step2. 在`proto/user`下创建文件`user.proto`

这一步就明确了有哪些服务是暴露给其他服务使用的

`user.proto`代码如下:

```proto
syntax = "proto3";

// 包名
package go.micro.service.user;

// 定义服务
service User {
  // 注册
  rpc Register(UserRegisterRequest) returns(UserRegisterResponse) {}

  // 登录
  rpc Login(UserLoginRequest) returns(UserLoginResponse) {}

  // 查询用户信息
  rpc GetUserInfo(UserInfoRequest) returns(UserInfoResponse) {}
}

// 查询用户信息请求体
message UserInfoRequest {
  // 账号(唯一)
  string userName = 1;
}

// 查询用户信息响应体
message UserInfoResponse {
  // 用户Id
  int64 userId = 1;

  // 账号(唯一)
  string userName = 2;

  // 用户名
  string firstName = 3;
}

// 用户注册请求体
message UserRegisterRequest {
  // 用户注册的账号
  string userName = 1;

  // 用户名
  string firstName = 2;

  // 用户密码
  string password = 3;
}

// 用户注册响应体
message UserRegisterResponse {
  // 表示注册是否成功的标量
  bool isSuccess = 1;
}

// 用户登录请求体
message UserLoginRequest {
  // 账号名
  string userName = 1;

  // 密码
  string password = 2;
}

// 用户登录响应体
message UserLoginResponse {
  // 表示登录是否成功的标量
  bool isSuccess = 1;
}
```

### 2.3 根据proto生成go代码

在工程的根目录下执行:

```
docker run --rm -v $(PWD):$(PWD) -w $(PWD) -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --go_out=./ --micro_out=./ ./proto/user/user.proto
恭喜，恭喜命令执行成功！%  
```

执行后的`proto/user`目录结构如下:

```
tree ./proto 
./proto
└── user
    ├── user.pb.go	// 基础文件
    ├── user.pb.micro.go	// micro的文件
    └── user.proto

1 directory, 3 files
```