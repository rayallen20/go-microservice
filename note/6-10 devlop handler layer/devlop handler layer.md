# devlop handler layer

## PART1. 加载common包

命令行中执行`go get github.com/rayallen20/common`

注:此处我将common仓库放到了github上,因为从慕课网的仓库上go get,需要私钥,我配置了但是不对.

[ssh配置公钥私钥的方法](https://blog.csdn.net/weixin_42209753/article/details/128478647)

## PART2. 实现`cart.CartHandler`接口

### 2.1 实现`AddCart()`方法

```go
package handler

import (
	"context"
	"git.imooc.com/cap1573/cart/domain/model"
	"git.imooc.com/cap1573/cart/domain/service"
	cart "git.imooc.com/cap1573/cart/proto/cart"
	"github.com/rayallen20/common"
)

type Cart struct {
	CartDataService service.ICartDataService
}

// AddCart 创建Cart信息
func (c *Cart) AddCart(ctx context.Context, request *cart.CartInfo, response *cart.ResponseAdd) (err error) {
	cart := &model.Cart{}
	common.SwapTo(request, cart)
	response.CartId, err = c.CartDataService.AddCart(cart)
	return err
}
```

注:此处已经引入了common包了,注意import处的`github.com/rayallen20/common`

### 2.2 实现`CleanCart()`方法

```go
// CleanCart 清空购物车
func (c *Cart) CleanCart(ctx context.Context, request *cart.Clean, response *cart.Response) error {
	err := c.CartDataService.CleanCart(request.UserId)
	if err != nil {
		return err
	}

	response.Msg = "购物车清空成功"
	return nil
}
```

### 2.3 实现`Incr()`方法

```go
// Incr 向指定购物车中添加指定数量的商品
func (c *Cart) Incr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	err := c.CartDataService.IncrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}

	response.Msg = "购物车添加成功"
	return nil
}
```

### 2.4 实现`Decr()`方法

```go
// Decr 从指定购物车中减少指定数量的商品
func (c *Cart) Decr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	err := c.CartDataService.DecrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}

	response.Msg = "购物车减少成功"
	return nil
}
```

### 2.5 实现`DeleteItemByID()`方法

```go
// DeleteItemByID 根据ID删除Cart信息
func (c *Cart) DeleteItemByID(ctx context.Context, request *cart.CartID, response *cart.Response) error {
	err := c.CartDataService.DeleteCart(request.Id)
	if err != nil {
		return err
	}
	
	response.Msg = "删除购物车成功"
	return nil
}
```

### 2.6 实现`GetAll()`方法

```go
// GetAll 获取指定用户的购物车
func (c *Cart) GetAll(ctx context.Context, request *cart.CartFindAll, response *cart.CartAll) error {
	cartAll, err := c.CartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}

	for _, cartOrm := range cartAll {
		cartBiz := &cart.CartInfo{}
		err = common.SwapTo(cartOrm, cartBiz)
		if err != nil {
			return err
		}
		response.CartInfo = append(response.CartInfo, cartBiz)
	}

	return nil
}
```

## PART3. 完整代码

`cart/handler/cart.go`完整代码如下:

```go
package handler

import (
	"context"
	"git.imooc.com/cap1573/cart/domain/model"
	"git.imooc.com/cap1573/cart/domain/service"
	cart "git.imooc.com/cap1573/cart/proto/cart"
	"github.com/rayallen20/common"
)

type Cart struct {
	CartDataService service.ICartDataService
}

// AddCart 创建Cart信息
func (c *Cart) AddCart(ctx context.Context, request *cart.CartInfo, response *cart.ResponseAdd) (err error) {
	cartOrm := &model.Cart{}
	err = common.SwapTo(request, cartOrm)
	if err != nil {
		return err
	}
	response.CartId, err = c.CartDataService.AddCart(cartOrm)
	return err
}

// CleanCart 清空购物车
func (c *Cart) CleanCart(ctx context.Context, request *cart.Clean, response *cart.Response) error {
	err := c.CartDataService.CleanCart(request.UserId)
	if err != nil {
		return err
	}

	response.Msg = "购物车清空成功"
	return nil
}

// Incr 向指定购物车中添加指定数量的商品
func (c *Cart) Incr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	err := c.CartDataService.IncrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}

	response.Msg = "购物车添加成功"
	return nil
}

// Decr 从指定购物车中减少指定数量的商品
func (c *Cart) Decr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	err := c.CartDataService.DecrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}

	response.Msg = "购物车减少成功"
	return nil
}

// DeleteItemByID 根据ID删除Cart信息
func (c *Cart) DeleteItemByID(ctx context.Context, request *cart.CartID, response *cart.Response) error {
	err := c.CartDataService.DeleteCart(request.Id)
	if err != nil {
		return err
	}

	response.Msg = "删除购物车成功"
	return nil
}

// GetAll 获取指定用户的购物车
func (c *Cart) GetAll(ctx context.Context, request *cart.CartFindAll, response *cart.CartAll) error {
	cartAll, err := c.CartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}

	for _, cartOrm := range cartAll {
		cartBiz := &cart.CartInfo{}
		err = common.SwapTo(cartOrm, cartBiz)
		if err != nil {
			return err
		}
		response.CartInfo = append(response.CartInfo, cartBiz)
	}

	return nil
}
```