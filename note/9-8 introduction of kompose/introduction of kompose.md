# introduction of kompose

## PART1. Kompose介绍

是一个将docker-compose的yaml文件快速转换为k8s能够部署文件的工具

## PART2. Kompose使用条件

- 必须具备一个已经安装好的k8s集群
- kubectl CLI工具必须能够连接到搭建好的K8S集群上

## PART3. Kompose的安装

Linux: `curl -l https://github.com/kubernetes/kompose/releases/download/v1.22.0/kompose-linux-amd64 -o kompose`

MacOS: `curl -l https://github.com/kubernetes/kompose/releases/download/v1.22.0/kompose-darwin-amd64 -o kompose`

`chmod +x kompose`

`sudo mv ./kompose /usr/local/bin/kompose`

## PART4. Kompose使用命令

- `kompose -f xxx.yml covert`:转化docker-compose文件
- `kubectl apply -f xxx.yaml`:从转化的文件中部署资源文件

