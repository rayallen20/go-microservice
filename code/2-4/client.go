package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	imooc "roach-imooc/proto/roach"
)

func main() {
	// 创建新的服务
	service := micro.NewService(
		micro.Name("roach.imooc.server"),
	)

	// 初始化
	service.Init()

	roachImooc := imooc.NewRoachService("roach.imooc.server", service.Client())

	res, err := roachImooc.SayHello(context.TODO(), &imooc.SayRequest{Message: "客户端发送的请求"})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", res.Answer)
}
