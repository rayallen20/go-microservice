package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	imooc "roach-imooc/proto/roach"
)

// RoachServer 需要实现接口RoachService
// 这个接口是protoc生成的 pb.micro.go
type RoachServer struct{}

func (r *RoachServer) SayHello(ctx context.Context, in *imooc.SayRequest, out *imooc.SayResponse) error {
	// 业务逻辑代码
	out.Answer = "SayHello的响应"
	return nil
}

func main() {
	// 创建新的服务
	service := micro.NewService(
		// 服务名 可以认为是服务的ID 是一个唯一标识符
		// Client通过这个name来找server
		micro.Name("roach.imooc.server"),
	)

	// 初始化服务
	service.Init()

	// 注册服务
	// 此处是pb自动生成的imooc.pb.go中的方法
	imooc.RegisterRoachHandler(service.Server(), new(RoachServer))

	// 运行服务
	err := service.Run()

	if err != nil {
		fmt.Println(err)
	}
}
