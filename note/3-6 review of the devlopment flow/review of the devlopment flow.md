# review of the devlopment flow

## PART1. 流程回顾

4个步骤:

1. 编写`proto`文件

	1.1 定义`service`
	
	1.2 定义`message`
	
	1.3 生成GO代码
	
	其中:生成的`.pb.go`文件中的接口`XXXService`,即为要实现的部分
	
2. `domain`层的开发

	2.1 定义ORM(`model`层)
	
	2.2 实现数据库操作(`repository`层)
	
	2.3 实现业务逻辑(`service`层,类似于我自己写代码时的`biz`层)
	
3. `handler`层的开发

	3.1 这一层要实现`.pb.go`文件中的接口`XXXService`
	
4. `main()`函数的开发

	4.1 创建服务
	
	4.2 服务初始化
	
	4.3 创建服务实例
	
	4.4 注册`handler`
	
	4.5 运行服务
	
## PART2. 编译

这一步如果有问题,就把`go.mod`中`require`的部分删除,然后再重新执行`go mod tidy`即可

```
2022-10-25 11:45:00  file=v2@v2.9.1/service.go:200 level=info Starting [service] go.micro.service.user
2022-10-25 11:45:00  file=grpc/grpc.go:864 level=info Server [grpc] Listening on [::]:60315
2022-10-25 11:45:00  file=grpc/grpc.go:697 level=info Registry [mdns] Registering node: go.micro.service.user-49fc48b9-a1be-44f5-ac67-c39aea46715e
```

可以看到,编译后的运行符合预期

注:

`mdns`:在没有注册中心时,通过dns实现服务发现和服务注册