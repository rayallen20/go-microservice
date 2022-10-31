# devlop the service which externally exposed

## PART1. 编写proto文件,确认对外暴露的服务

```proto
syntax = "proto3";

package go.micro.service.category;

service Category {
	// CreateCategory 创建和更新的时候都会用到
	rpc CreateCategory(CategoryRequest) returns (CreateCategoryResponse) {}

	// UpdateCategory 更新品类信息
	rpc UpdateCategory(CategoryRequest) returns (UpdateCategoryResponse) {}

	// DeleteCategory 删除品类信息
	rpc DeleteCategory(DeleteCategoryRequest) returns (DeleteCategoryResponse) {}

	// FindCategoryByName 根据名称查询品类信息
	rpc FindCategoryByName(FindByNameRequest) returns (CategoryResponse) {}

	// FindCategoryById 根据id查询品类信息
	rpc FindCategoryById(FindByIdRequest) returns (CategoryResponse) {}

	// FindCategoryByLevel 根据层级查询品类信息
	rpc FindCategoryByLevel(FindByLevelRequest) returns (FindAllResponse) {}

	// FindCategoryByParent 根据父品类查询品类信息
	rpc FindCategoryByParent(FindByParentRequest) returns (FindAllResponse) {}

	// FindAllCategory 查询所有品类信息
	rpc FindAllCategory(FindAllRequest) returns (FindAllResponse) {}
}

// CategoryRequest 创建品类和删除品类的请求message
message CategoryRequest {
	// category_name 品类名称
	string category_name = 1;

	// category_level 品类层级 用于多级分类
	uint32 category_level = 2;

	// 父级品类的id
	int64 category_parent = 3;

	// 品类图片
	string category_image = 4;

	// 品类描述
	string category_description = 5;
}

// CreateCategoryResponse 创建品类的响应
message CreateCategoryResponse {
	string message = 1;

	// 品类id
	int64 category_id = 2;
}

// UpdateCategoryResponse 更新品类的响应
message UpdateCategoryResponse {
	string message = 1;
}

// DeleteCategoryRequest 删除品类的请求
message DeleteCategoryRequest {
	// 待删除的品类id
	int64 category_id = 1;
}

// DeleteCategoryResponse 删除品类的响应
message DeleteCategoryResponse {
	string message = 1;
}

// FindByNameRequest 根据名称查找品类的请求
message FindByNameRequest {
	string category_name = 1;
}

// CategoryResponse 品类响应
message CategoryResponse {
	// 品类id
	int64 id = 1;

	// 品类名称
	string category_name = 2;

	// 品类层级
	uint32 category_level = 3;

	// 品类的父品类id
	int64 category_parent = 4;

	// 品类图片
	string category_image = 5;

	// 品类描述
	string category_description = 6;
}

// FindByIdRequest 根据id查找品类的请求
message FindByIdRequest {
	int64 category_id = 1;
}

// FindByLevelRequest 根据层级查询品类信息请求
message FindByLevelRequest {
	uint32 level = 1;
}

// FindByParentRequest 根据父品类查询品类请求
message FindByParentRequest {
	int64 parent_id = 1;
}

// FindAllRequest 查询所有品类信息请求
message FindAllRequest {
	// 可以写个分页啥的 此处演示 就不写了
}

// FindAllResponse 查询所有品类信息响应
message FindAllResponse {
	// repeated相当于list
	repeated CategoryResponse category = 1;
}
```

## PART2. 生成GO代码

此处使用了Makefile简化命令,Makefile内容如下:

```Makefile


.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/category/category.proto

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o category-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t category-service:latest
```

生成代码:

```
make proto
sudo docker run --rm -v /Users/yanglei/Desktop/user/git.imooc.com/cap1573/category:/Users/yanglei/Desktop/user/git.imooc.com/cap1573/category -w /Users/yanglei/Desktop/user/git.imooc.com/cap1573/category -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/category/category.proto
Password:
恭喜，恭喜命令执行成功！%  
```

```
tree ./proto/category 
./proto/category
├── category.pb.go
├── category.pb.micro.go
└── category.proto

0 directories, 3 files
```