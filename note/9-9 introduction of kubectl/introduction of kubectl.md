# introduction of kubectl

## PART1. kubectl常用命令分类

- 命令式资源管理
- 资源查看
- 容器管理

## PART2. 命令式资源管理

- 创建:
	
	- `create`:创建一个资源
	- `expose`:暴露一个资源

- 更新:

	- `scale`:扩展资源
	- `annotate`:添加备注
	- `label`:标签

- 删除:

	- `delete`:删除资源

## PART3. 资源查看

- `get`:最常用的查看命令,显示一个或多个资源的详细信息
- `describe`:`describe`命令同样用于查看资源信息,但`get`只输出资源本身的信息,`describe`聚合了相关资源的信息并输出

## PART4. 容器管理

- `log`:查看容器log
- `exec`:执行命令
- `cp`:用于在容器与物理机文件的拷贝

## PART5. kubectl语法

语法:`kubectl [command] [TYPE] [NAME] [flags]`

- `get`:`kubectl get po,svc`
- `log`:`kubectl log mysql`
- `exec`:`kubectl exec -it mysql /bin/bash`

## PART6. K8S中常用资源类型的简写

- ing:ingress
- no:集群中的nodes节点
- ns:namespace
- rs:replica sets
- svc:service
- po:pod