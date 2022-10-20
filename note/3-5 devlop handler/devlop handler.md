# devlop handler

## PART1. 实现handler层

handler层需要实现protoc生成的`user.micro.go`中的`UserHandler`接口.也就是实现对外暴露的服务

实际上handler层连接了service层(`domain/service`)和proto层(`proto/user`)

`handler/user.go`:

```go
package handler

import (
	"context"
	"git.imooc.com/rayallen20c/user/domain/model"
	"git.imooc.com/rayallen20c/user/domain/service"
	user "git.imooc.com/rayallen20c/user/proto/user"
)

// User 需实现protoc生成的user.pb.micro.go中的UserHandler接口
type User struct {
	UserDataService service.IUserDataService
}

// Register 用户注册
func (u *User) Register(ctx context.Context, in *user.UserRegisterRequest, out *user.UserRegisterResponse) error {
	userRegister := &model.User{
		UserName:     in.UserName,
		FirstName:    in.FirstName,
		HashPassword: in.Password,
	}

	_, err := u.UserDataService.AddUser(userRegister)
	if err != nil {
		return err
	}
	out.IsSuccess = true
	return nil
}

// Login 用户登录
func (u *User) Login(ctx context.Context, in *user.UserLoginRequest, out *user.UserLoginResponse) error {
	isOk, err := u.UserDataService.CheckPwd(in.UserName, in.Password)
	if err != nil {
		return err
	}

	if !isOk {
		return err
	}

	out.IsSuccess = true
	return nil
}

// GetUserInfo 查询用户信息
func (u *User) GetUserInfo(ctx context.Context, in *user.UserInfoRequest, out *user.UserInfoResponse) error {
	userInfo, err := u.UserDataService.FindUserByName(in.UserName)
	if err != nil {
		return err
	}

	out = UserInfoToResponse(userInfo)
	return nil
}

// UserInfoToResponse 将user ORM转换为响应对象
func UserInfoToResponse(userInfo *model.User) (resp *user.UserInfoResponse) {
	resp = &user.UserInfoResponse{
		UserId:    userInfo.Id,
		UserName:  userInfo.UserName,
		FirstName: userInfo.FirstName,
	}
	return resp
}
```

## PART2. 修改main()函数

```go
package main

import (
	"fmt"
	"git.imooc.com/rayallen20c/user/domain/repository"
	service2 "git.imooc.com/rayallen20c/user/domain/service"
	"git.imooc.com/rayallen20c/user/handler"
	user "git.imooc.com/rayallen20c/user/proto/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
)

func main() {
	// 服务参数设置
	srv := micro.NewService(
		// 服务名称 通常是和proto文件中的package名相同的
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
	)

	// 服务初始化
	srv.Init()

	// DB初始化
	db, err := gorm.Open("mysql", "root:123456@/micro?charset=utf8&parseTime=True&loc=local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// 禁止复表名称
	db.SingularTable(true)

	// 这段代码只执行1次
	rp := repository.NewIUserRepository(db)
	rp.InitTable()

	// 创建服务实例
	userDataService := service2.NewUserDataService(rp)

	// 注册handler
	err = user.RegisterUserHandler(srv.Server(), &handler.User{UserDataService: userDataService})

	// 运行服务
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
```

1. `go-micro`要使用v2版本(`go get github.com/micro/go-micro/v2`然后再`go mod tidy`即可)
2. DB的连接句柄和表初始化
3. 绑定handler层的User对象(`handler.User`)到服务
4. 运行服务