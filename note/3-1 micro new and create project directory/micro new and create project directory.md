# micro new and create project directory

## PART1. docker中使用micro & 项目目录搭建

### 1.1 项目目录搭建

#### 1.1.1 使用micro new生成项目初始目录

- step1. 拉取镜像

```
docker pull micro/micro
```

- step2. 生成项目代码

```
sudo docker run --rm -v $(pwd):$(pwd) -w $(pwd) \
> micro/micro new user
Password:
Creating service user

.
├── micro.mu
├── main.go
├── generate.go		// 模板文件生成的 开发时很少用
├── handler	// 定义暴露出的服务 可以认为是controller 
│   └── user.go
├── proto
│   └── user.proto
├── Dockerfile
├── Makefile
├── README.md
├── .gitignore
└── go.mod


download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases

compile the proto file user.proto:

cd user
make init
go mod vendor
make proto
```

生成后项目目录结构如下:

```
tree ./     
./
├── Dockerfile
├── Makefile
├── README.md
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── micro.mu
└── proto
    └── user.proto

2 directories, 10 files
```

#### 1.1.2 在项目初始目录上添加目录

- step1. 新建目录domain

可以认为该目录是处理业务逻辑的目录

创建后目录结构如下:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── domain
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── micro.mu
└── proto
    └── user.proto

3 directories, 10 files
```

- step2. 在domain目录下新建目录model

该目录用于存放数据模型(ORM)

创建后目录结构如下:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── domain
│   └── model
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── micro.mu
└── proto
    └── user.proto

4 directories, 10 files
```

- step3. 在domain目录下新建目录repository

该目录用于操作数据库

创建后目录结构如下:

```
tree ./
./
├── Dockerfile
├── Makefile
├── README.md
├── domain
│   ├── model
│   └── repository
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── micro.mu
└── proto
    └── user.proto

5 directories, 10 files
```

- step4. 在domain目录下新建目录service

该目录指代的service不是micro service中的service,而是业务逻辑.该目录主要和repository和model交互

创建后目录结构如下:

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
    └── user.proto

6 directories, 10 files
```