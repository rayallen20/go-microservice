# devlop proto which in product domain

## PART1. 生成代码

```
docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=7DD47DEF3E0D096A cap1573/cap-micro new git.imooc.com/cap1573/product
Creating service go.micro.service.product in git.imooc.com/cap1573/product

.
├── main.go
├── generate.go
├── plugin.go
├── handler
│   └── product.go
├── domain/model
│   └── product.go
├── domain/repository
│   └── product_repository.go
├── domain/service
│   └── product_data_service.go
├── common
│   └── README.md
├── proto/product
│   └── product.proto
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

compile the proto file product.proto:

cd git.imooc.com/cap1573/product
protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/product/product.proto
```

## PART2. 编写proto

`product/proto/product/product.proto`:

```proto
syntax = "proto3";

package go.micro.service.product;

service Product {
	// AddProduct 新增商品
	rpc AddProduct(ProductInfo) returns (ResponseProduct) {}

	// FindProductByID 根据ID查询商品信息
	rpc FindProductByID(RequestID) returns (ProductInfo) {}

	// UpdateProduct 更新商品信息
	rpc UpdateProduct(ProductInfo) returns (Response) {}

	// DeleteProductByID 删除商品信息
	rpc DeleteProductByID(RequestID) returns (Response) {}

	// FindAllProduct 查询全部商品信息
	rpc FindAllProduct(RequestAll) returns (AllProduct) {}
}

// ProductInfo 商品信息
message ProductInfo {
	// id 商品id
	int64 id = 1;

	// product_name 商品名称
	string product_name = 2;

	// product_sku 商品款名 可以认为类似商品的条形码 是商品的唯一标识
	string product_sku = 3;

	// product_price 商品价格
	double product_price = 4;

	// product_description 商品描述
	string product_description = 5;

	// product_category_id 商品所属品类id
	int64 product_category_id = 6;

	// product_image 商品图片集合
	repeated ProductImage product_image = 7;

	// product_size 商品尺码集合
	repeated ProductSize product_size = 8;

	// ProductSeo 商品的SEO优化信息
	ProductSeo product_seo = 9;
}

// ProductImage 商品图片
message ProductImage {
	// id 图片id
	int64 id = 1;

	// image_name 图片名称
	string image_name = 2;

	// image_code 图片编码 用于保证对图片操作的幂等性使用 可以认为是图片的唯一标识
	string image_code = 3;

	// image_url 图片URL
	string image_url = 4;
}

// ProductSize 商品尺码
message ProductSize {
	// id 尺码id
	int64 id = 1;

	// size_name 尺码名称
	string size_name = 2;

	// size_code 尺码编码 同样是用于保证操作幂等性的 可以认为是尺码的唯一标识
	string size_code = 3;
}

// ProductSeo 商品的SEO优化信息
message ProductSeo {
	// id seo优化信息id
	int64 id = 1;

	// seo_title 展示给搜索引擎的标题
	string seo_title = 2;

	// seo_keyword 展示给搜索引擎的关键词
	string seo_keyword = 3;

	// seo_description 展示给搜索引擎的描述信息
	string seo_description = 4;

	// seo_code seo信息编码 同样是用于保证操作幂等性的 可以认为是seo信息的唯一标识
	string seo_code = 5;
}

// ResponseProduct 添加商品的响应信息
message ResponseProduct {
	// product_id 商品信息id
	int64 product_id = 1;
}

// RequestID 更新/删除商品接口的请求体
message RequestID {
	// product_id 商品信息id
	int64  product_id = 1;
}

// Response 更新/删除商品接口的响应体
message Response {
	// msg 类似于响应码 表示操作结果的唯一标识
	string msg = 1;
}

// RequestAll 查询全部商品信息接口的请求体 实际上无任何参数
message RequestAll {

}

// AllProduct 查询全部商品信息接口的响应体
message AllProduct {
	// product_info 商品信息集合
	repeated ProductInfo product_info = 1;
}
```

## PART3. 生成代码

Makefile:

```Makefile


.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/product/product.proto

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o product-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t product-service:latest
```

在项目的根目录下执行:

```
make proto
sudo docker run --rm -v /product:/product -w /product -e ICODE=7DD47DEF3E0D096A cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/product/product.proto
Password:
恭喜，恭喜命令执行成功！
```

注:实际上就是执行的Makefile中的proto部分的命令

```
tree ./proto/product 
./proto/product
├── product.pb.go
├── product.pb.micro.go
└── product.proto

0 directories, 3 files
```





















