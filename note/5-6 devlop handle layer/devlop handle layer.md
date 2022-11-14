# devlop handle layer

`handle.Product`需要实现`product.pb.micro.go`中的接口`ProductHandler`.

注:开发前需删除`product/handler/product.go`中除结构体定义外的所有方法.

## PART1. 实现`AddProduct()`方法

```go
// AddProduct 新增商品
func (p *Product) AddProduct(ctx context.Context, request *product.ProductInfo, response *product.ResponseProduct) error {
	productOrm := &model.Product{}
	err := common.SwapTo(request, productOrm)
	if err != nil {
		return err
	}
	
	productID, err := p.ProductDataService.AddProduct(productOrm)
	if err != nil {
		return err
	}
	
	response.ProductId = productID
	return nil
}
```

注:此处的`common.SwapTo()`函数和上一案例(category)中的`common.SwapTo()`完全相同,故代码不再粘贴一遍了.

## PART2. 实现`FindProductByID()`方法

```go
// FindProductByID 根据ID查询商品信息
func (p *Product) FindProductByID(ctx context.Context, request *product.RequestID, response *product.ProductInfo) error {
	productOrm, err := p.ProductDataService.FindProductByID(request.ProductId)
	if err != nil {
		return err
	}

	err = common.SwapTo(productOrm, response)
	if err != nil {
		return err
	}
	return nil
}
```

## PART3. 实现`UpdateProduct()`方法

```go
// UpdateProduct 更新商品信息
func (p *Product) UpdateProduct(ctx context.Context, request *product.ProductInfo, response *product.Response) error {
	productOrm := &model.Product{}
	err := common.SwapTo(request, productOrm)
	if err != nil {
		return err
	}
	
	err = p.ProductDataService.UpdateProduct(productOrm)
	if err != nil {
		return err
	}
	
	response.Msg = "更新成功"
	return nil
}
```

## PART4. 实现`DeleteProductByID()`方法

```go
// DeleteProductByID 删除商品信息
func (p *Product) DeleteProductByID(ctx context.Context, request *product.RequestID, response *product.Response) error {
	err := p.ProductDataService.DeleteProduct(request.ProductId)
	if err != nil {
		return err
	}
	
	response.Msg = "删除成功"
	return nil
}
```

## PART5. 实现`FindAllProduct()`方法

```go
// FindAllProduct 查询全部商品信息
func (p *Product) FindAllProduct(ctx context.Context, request *product.RequestAll, response *product.AllProduct) error {
	productOrmCollection, err := p.ProductDataService.FindAllProduct()
	if err != nil {
		return err
	}
	
	for _, productOrm := range productOrmCollection {
		productResp := &product.ProductInfo{}
		err = common.SwapTo(productOrm, productResp)
		if err != nil {
			return err
		}
		
		response.ProductInfo = append(response.ProductInfo, productResp)
	}
	
	return nil
}
```

## PART6. 完整handler层代码

`product/handler/product.go`:

```go
package handler

import (
	"context"
	"git.imooc.com/cap1573/product/common"
	"git.imooc.com/cap1573/product/domain/model"
	"git.imooc.com/cap1573/product/domain/service"
	product "git.imooc.com/cap1573/product/proto/product"
)

type Product struct {
	ProductDataService service.IProductDataService
}

// AddProduct 新增商品
func (p *Product) AddProduct(ctx context.Context, request *product.ProductInfo, response *product.ResponseProduct) error {
	productOrm := &model.Product{}
	err := common.SwapTo(request, productOrm)
	if err != nil {
		return err
	}

	productID, err := p.ProductDataService.AddProduct(productOrm)
	if err != nil {
		return err
	}

	response.ProductId = productID
	return nil
}

// FindProductByID 根据ID查询商品信息
func (p *Product) FindProductByID(ctx context.Context, request *product.RequestID, response *product.ProductInfo) error {
	productOrm, err := p.ProductDataService.FindProductByID(request.ProductId)
	if err != nil {
		return err
	}

	err = common.SwapTo(productOrm, response)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProduct 更新商品信息
func (p *Product) UpdateProduct(ctx context.Context, request *product.ProductInfo, response *product.Response) error {
	productOrm := &model.Product{}
	err := common.SwapTo(request, productOrm)
	if err != nil {
		return err
	}

	err = p.ProductDataService.UpdateProduct(productOrm)
	if err != nil {
		return err
	}

	response.Msg = "更新成功"
	return nil
}

// DeleteProductByID 删除商品信息
func (p *Product) DeleteProductByID(ctx context.Context, request *product.RequestID, response *product.Response) error {
	err := p.ProductDataService.DeleteProduct(request.ProductId)
	if err != nil {
		return err
	}

	response.Msg = "删除成功"
	return nil
}

// FindAllProduct 查询全部商品信息
func (p *Product) FindAllProduct(ctx context.Context, request *product.RequestAll, response *product.AllProduct) error {
	productOrmCollection, err := p.ProductDataService.FindAllProduct()
	if err != nil {
		return err
	}

	for _, productOrm := range productOrmCollection {
		productResp := &product.ProductInfo{}
		err = common.SwapTo(productOrm, productResp)
		if err != nil {
			return err
		}

		response.ProductInfo = append(response.ProductInfo, productResp)
	}

	return nil
}
```