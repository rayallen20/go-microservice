# generate code of category module

## PART1. 电商微服务分类模块开发

### 1.1 分类模块目录建立

此处直接使用上课现成的模板,这部分后续要考虑如何手动实现.甚至考虑如何打一个自己的go-micro镜像出来用于研发.

```
pwd  
/Users/yanglei/Desktop/user
```

```
docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=7DD47DEF3E0D096A cap1573/cap-micro new git.imooc.com/cap1573/category
Creating service go.micro.service.category in git.imooc.com/cap1573/category

.
├── main.go
├── generate.go
├── plugin.go
├── handler
│   └── category.go
├── domain/model
│   └── category.go
├── domain/repository
│   └── category_repository.go
├── domain/service
│   └── category_data_service.go
├── common
│   └── README.md
├── proto/category
│   └── category.proto
├── Dockerfile
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases

download protobuf for micro:

go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/micro/micro/v2/cmd/protoc-gen-micro

compile the proto file category.proto:

cd git.imooc.com/cap1573/category
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/category/category.proto
```

```
tree ./git.imooc.com     
./git.imooc.com
└── cap1573
    └── category
        ├── Dockerfile
        ├── Makefile
        ├── README.md
        ├── common	// 存放公共方法
        │   └── README.md
        ├── domain
        │   ├── model	// ORM
        │   │   └── category.go
        │   ├── repository		// 对ORM的CRUD
        │   │   └── category_repository.go
        │   └── service		// CRUD的逻辑部分
        │       └── category_data_service.go
        ├── generate.go
        ├── go.mod
        ├── handler
        │   └── category.go
        ├── main.go
        ├── plugin.go
        └── proto
            └── category
                └── category.proto

10 directories, 13 files
```