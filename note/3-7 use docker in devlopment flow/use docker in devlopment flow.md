# use docker in devlopment flow

## PART1. 编写Dockerfile

Dockerfile:用于构建镜像的文本文件.文本内容包含了构建镜像所需的指令和说明

### 1.1 Dockerfile常用指令

- `FROM`: 定制的镜像都是基于`FROM`的镜像,后续的操作都是基于这个镜像
- `RUN`: 用于执行后边的命令
- `COPY/ADD`:复制指令,若目标为压缩文件,则`ADD`指令会解压,而`COPY`指令不会
- `CMD/ENTRYPOINT`:用于启动程序
- `ENV`:设置环境变量
- `EXPOSE`:声明暴露的端口,但仅仅是声明,不会监听
- `WORKDIR`:指定工作目录
- `USER`:用于指定后续命令的用户和用户组

### 1.2 使用Dockerfile构建微服务

#### 1.2.1 编译

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user *.go
```

```
tree ./ -L 1
./
├── Dockerfile
├── Makefile
├── README.md
├── domain
├── generate.go
├── go.mod
├── go.sum
├── handler
├── main.go
├── micro.mu
├── proto
└── user

3 directories, 9 files
```

可以看到,多了一个二进制文件`user`

#### 1.2.2 编写Dockerfile

Dockerfile内容如下:

```Dockerfile
FROM alpine
ADD user /user
ENTRYPOINT [ "/user" ]
```

#### 1.2.3 构建镜像

```
docker build -t user:v1.0 -f Dockerfile .   
[+] Building 1.2s (7/7) FINISHED                                                                                                                                                                                                                                         
 => [internal] load build definition from Dockerfile                                                                                                                                                                                                                0.0s
 => => transferring dockerfile: 92B                                                                                                                                                                                                                                 0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                                   0.0s
 => => transferring context: 2B                                                                                                                                                                                                                                     0.0s
 => [internal] load metadata for docker.io/library/alpine:latest                                                                                                                                                                                                    0.0s
 => [internal] load build context                                                                                                                                                                                                                                   0.6s
 => => transferring context: 31.92MB                                                                                                                                                                                                                                0.6s
 => CACHED [1/2] FROM docker.io/library/alpine                                                                                                                                                                                                                      0.0s
 => [2/2] ADD user /user                                                                                                                                                                                                                                            0.3s
 => exporting to image                                                                                                                                                                                                                                              0.2s
 => => exporting layers                                                                                                                                                                                                                                             0.2s
 => => writing image sha256:4cb6442a66cca9eb4f93c26fd3e1a8d4a2dc99550d652fa02d6524357f354eb0                                                                                                                                                                        0.0s
 => => naming to docker.io/library/user:v1.0                                                                                                                                                                                                                        0.0s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
```

## PART2. 编写Makefile(非必须)

### 2.1 Makefile的作用

简化命令.在本例中,可以简化对`protoc`、`go build`、`docker build`指令的使用.简单理解就是将这些指令,通过在`Makefile`中写一些别名的方式,简化这些指令

### 2.2 编写Makefile

```Makefile

GOPATH:=$(shell go env GOPATH)
.PHONY: init

.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --go_out=./ --micro_out=./ ./proto/user/user.proto
	
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: dockerBuild
dockerBuild:
	sudo docker build -t user:v1.0 -f Dockerfile .

```

### 2.3 使用Makefile

```
make proto
sudo docker run --rm -v /Users/yanglei/Desktop/user/user:/Users/yanglei/Desktop/user/user -w /Users/yanglei/Desktop/user/user -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --go_out=./ --micro_out=./ ./proto/user/user.proto
恭喜，恭喜命令执行成功！% 
```