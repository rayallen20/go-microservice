# introduce of docker-compose

## PART1. docker-compose介绍

- 用于定义和运行多容器Docker应用程序的工具
- 使用YML文件来配置应用程序需要的所有服务
- 使用一个命令,可以创建并启动所有服务

## PART2. docker-compose的使用

- step1. 使用Dockerfile定义应用程序的环境
- step2. 使用docker-compose.yml定义构成应用程序的服务
- step3. 执行`docker-compose up`启动即可

## PART3. docker-compose的安装

### 3.1 Linux安装docker-compose

- step1. 下载

```
sudo curl -L "https://github.com/docker/compose/releases/download/1.28.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

- step2. 提权

```
sudo chmod +x /usr/local/bin/docker-compose
```

- step3. 建立软链接

```
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```

### 3.2 MacOS/Windows

Docker Toolbox中已包含

检测:

```
docker-compose --version
Docker Compose version v2.12.0
```

## PART4. docker-compose的yml文件

- 通常使用docker-compose.yml定义构成应用程序的服务
- 用户通过模板文件来定义一组相关联的应用容器
- 标准模板文件应该包含version、services、networks三大部分,最关键的是services部分

## PART5. docker-compose常用命令

- `docker-compose up -d`:后台运行
- `docker-compose ps`:列出项目中所有的容器
- `docker-compose down`:停止和删除容器与网络

