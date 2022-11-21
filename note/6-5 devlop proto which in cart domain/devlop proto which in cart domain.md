# devlop proto which in cart domain

## PART1. 创建工程目录

```
docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=3DCF426135DEE96B cap1573/cap-micro new git.imooc.com/cap1573/cart
```

```
tree ./ -L 1
./
├── Dockerfile
├── Makefile
├── README.md
├── common
├── domain
├── generate.go
├── go.mod
├── handler
├── main.go
├── plugin.go
└── proto

4 directories, 7 files

```

注:此处生成的是服务端,API层是其客户端,对其进行调用.

## PART2. 编写proto

`cart/proto/cart.proto`:

```proto
syntax = "proto3";

package go.micro.service.cart;

service Cart {
	// AddCart 添加商品到购物车
	rpc AddCart(CartInfo) returns (ResponseAdd) {}

	// CleanCart 清空购物车
	rpc CleanCart(Clean) returns (Response) {}

	// Incr 添加购物车中商品的数量
	rpc Incr(Item) returns (Response) {}

	// Decr 减少购物车中商品的数量
	rpc Decr(Item) returns (Response) {}

	// DeleteItemByID 删除购物车
	rpc DeleteItemByID(cartID) returns (Response) {}

	// GetAll 查看所有购物车信息
	rpc GetAll(CartFindAll) returns (CartAll) {}
}

// CartInfo 购物车信息
message CartInfo {
	// id 购物车id
	int64 id = 1;

	// user_id 用户id
	int64 user_id = 2;

	// product_id 商品id
	int64 product_id = 3;

	// size_id 商品尺寸id
	int64 size_id = 4;

	// num 商品添加的件数
	int64 num = 5;
}

// ResponseAdd 添加商品到购物车的响应
message ResponseAdd {
	// cart_id 购物车id
	int64 cart_id = 1;

	// msg 添加操作的结果信息
	string msg = 2;
}

// Clean 清空购物车请求参数
message Clean {
	// user_id 用户id
	int64 user_id = 1;
}

// Response 操作结果的响应
message Response {
	// msg 表示操作结果的信息
	string msg = 1;
}

// Item 添加/减少购物车中商品的请求参数
message Item {
	// id 购物车id
	int64 id = 1;

	// change_num 变更后的数量
	int64 change_num = 2;
}

// cartID 删除购物车操作的请求参数
message cartID {
	// id 购物车id
	int64 id = 1;
}

// CartFindAll 查找所有购物车操作的请求参数
message CartFindAll {
	// user_id 用户id
	int64 user_id = 1;
}

// CartAll 查找所有购物车操作的响应参数
message CartAll {
	// cart_info 购物车信息集合
	repeated CartInfo cart_info = 1;
}
```

## PART3. 生成代码

```
 pwd
/cart
docker run --rm -v $(pwd):$(pwd) -w $(pwd) -e ICODE=3DCF426135DEE96B cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/cart/cart.proto
恭喜，恭喜命令执行成功！
```

```
tree ./proto 
./proto
└── cart
    ├── cart.pb.go
    ├── cart.pb.micro.go
    └── cart.proto

1 directory, 3 files
```