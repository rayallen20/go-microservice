# devlop handler

## PART1. 实现公共函数

本部分代码均在`user/git.imooc.com/cap1573/category/common/`下编写

### 1.1 实现从请求结构体中获取请求参数并赋值给ORM的函数`SwapTo`

- step1. 创建文件`swap.go`

```
tree ./common -L 1
./common
├── README.md
└── swap.go

0 directories, 2 files
```

- step2. 编写代码

```go
package common

import "encoding/json"

// SwapTo 通过json tag将请求结构体中的字段值赋值给Category结构体(本例中由于没有Biz层 故使用ORM替代)
// 也就是说实际上这个category应该是Service层的结构体 且本例中因为ORM结构体中每个字段的json tag
// 和请求结构体中对应字段的json tag相同 所以可以用这种抽象的方式
func SwapTo(request, category interface{}) (err error) {
	dataByte, err := json.Marshal(request)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataByte, category)
}
```

## PART2. 实现handler层

和之前的user一样,handler层要实现生成的`.pb.micro.go`中的接口`XXXHandler`

注:在开始编写后续代码之前,将`user/git.imooc.com/cap1573/category/handler/category.go`中,除了结构体定义之外部分的所有代码都删掉,如下:

```go
package handler

import (
	"git.imooc.com/cap1573/category/domain/service"
)

type Category struct {
	CategoryDataService service.ICategoryDataService
}
```

### 2.1 实现创建品类的方法

```go
package handler

import (
	"context"
	"git.imooc.com/cap1573/category/common"
	"git.imooc.com/cap1573/category/domain/model"
	"git.imooc.com/cap1573/category/domain/service"
	category "git.imooc.com/cap1573/category/proto/category"
)

type Category struct {
	CategoryDataService service.ICategoryDataService
}

// CreateCategory 提供创建分类的服务
func (c *Category) CreateCategory(ctx context.Context, request *category.CategoryRequest, response *category.CreateCategoryResponse) error {
	categoryOrm := &model.Category{}
	// 从请求结构体中直接将请求中的参数赋值给ORM
	err := common.SwapTo(request, categoryOrm)
	if err != nil {
		return err
	}
	categoryId, err := c.CategoryDataService.AddCategory(categoryOrm)
	if err != nil {
		return err
	}

	// 返回的信息(或者说结构体)也是作为入参传到方法中来的
	// 这一点和平常写代码不太一样
	response.CategoryId = categoryId
	response.Message = "分类添加成功"
	return nil
}
```

### 2.2 实现更新品类的方法

```go
// UpdateCategory 更新品类信息
func (c *Category) UpdateCategory(ctx context.Context, request *category.CategoryRequest, response *category.UpdateCategoryResponse) error {
	categoryOrm := &model.Category{}
	err := common.SwapTo(request, categoryOrm)
	if err != nil {
		return err
	}

	err = c.CategoryDataService.UpdateCategory(categoryOrm)
	if err != nil {
		return err
	}
	response.Message = "分类更新成功"
	return nil
}
```

### 2.3 实现删除品类的方法

```go
// DeleteCategory 删除品类信息
func (c *Category) DeleteCategory(ctx context.Context, request *category.DeleteCategoryRequest, response *category.DeleteCategoryResponse) error {
	err := c.CategoryDataService.DeleteCategory(request.CategoryId)
	if err != nil {
		return err
	}

	response.Message = "分类删除成功"
	return nil
}
```

### 2.4 实现根据名称查询品类信息的方法

```go
// FindCategoryByName 根据名称查询品类信息
func (c *Category) FindCategoryByName(ctx context.Context, request *category.FindByNameRequest, response *category.CategoryResponse) error {
	categoryOrm, err := c.CategoryDataService.FindCategoryByName(request.CategoryName)
	if err != nil {
		return err
	}

	err = common.SwapTo(categoryOrm, response)
	if err != nil {
		return err
	}

	return nil
}
```

### 2.5 实现根据id查询品类信息的方法

```go
// FindCategoryById 根据id查询品类信息
func (c *Category) FindCategoryById(ctx context.Context, request *category.FindByIdRequest, response *category.CategoryResponse) error {
	categoryOrm, err := c.CategoryDataService.FindCategoryByID(request.CategoryId)
	if err != nil {
		return err
	}

	err = common.SwapTo(categoryOrm, response)
	if err != nil {
		return err
	}

	return nil
}
```

### 2.6 实现根据层级查询品类信息的方法

#### 2.6.1 实现将ORM集合转换为1个FindAllResponse的响应的函数

```go
// categoriesToResponse 将ORM集合转换为1个FindAllResponse的响应
func categoriesToResponse(categories []model.Category, response *category.FindAllResponse) {
	for _, categoryOrm := range categories {
		categoryResponse := &category.CategoryResponse{}
		err := common.SwapTo(categoryOrm, categoryResponse)
		if err != nil {
			log.Error(err)
			break
		}
		response.Category = append(response.Category, categoryResponse)
	}
}
```

#### 2.6.2 实现根据层级查询品类信息的方法

```go
// FindCategoryByLevel 根据层级查询品类信息
func (c *Category) FindCategoryByLevel(ctx context.Context, request *category.FindByLevelRequest, response *category.FindAllResponse) error {
	categoryOrmCollection, err := c.CategoryDataService.FindCategoryByLevel(request.Level)
	if err != nil {
		return err
	}
	categoriesToResponse(categoryOrmCollection, response)
	return nil
}
```

### 2.7 实现根据父品类查询品类信息的方法

```go
// FindCategoryByParent 根据父品类查询品类信息
func (c *Category) FindCategoryByParent(ctx context.Context, request *category.FindByParentRequest, response *category.FindAllResponse) error {
	categoryOrmCollection, err := c.CategoryDataService.FindCategoryByParent(request.ParentId)
	if err != nil {
		return err
	}
	
	categoriesToResponse(categoryOrmCollection, response)
	return nil
}
```

### 2.8 实现查询所有品类信息的方法

```go
// FindAllCategory 查询所有品类信息
func (c *Category) FindAllCategory(ctx context.Context, request *category.FindAllRequest, response *category.FindAllResponse) error {
	categoryOrmCollection, err := c.CategoryDataService.FindAllCategory()
	if err != nil {
		return err
	}

	categoriesToResponse(categoryOrmCollection, response)
	return nil
}
```

### 2.9 handler层全部代码

`user/git.imooc.com/cap1573/category/handler/category.go`全部代码如下:

```go
package handler

import (
	"context"
	"git.imooc.com/cap1573/category/common"
	"git.imooc.com/cap1573/category/domain/model"
	"git.imooc.com/cap1573/category/domain/service"
	category "git.imooc.com/cap1573/category/proto/category"
	log "github.com/micro/go-micro/v2/logger"
)

type Category struct {
	CategoryDataService service.ICategoryDataService
}

// CreateCategory 提供创建分类的服务
func (c *Category) CreateCategory(ctx context.Context, request *category.CategoryRequest, response *category.CreateCategoryResponse) error {
	categoryOrm := &model.Category{}
	// 从请求结构体中直接将请求中的参数赋值给ORM
	err := common.SwapTo(request, categoryOrm)
	if err != nil {
		return err
	}
	categoryId, err := c.CategoryDataService.AddCategory(categoryOrm)
	if err != nil {
		return err
	}

	// 返回的信息(或者说结构体)也是作为入参传到方法中来的
	// 这一点和平常写代码不太一样
	response.CategoryId = categoryId
	response.Message = "分类添加成功"
	return nil
}

// UpdateCategory 更新品类信息
func (c *Category) UpdateCategory(ctx context.Context, request *category.CategoryRequest, response *category.UpdateCategoryResponse) error {
	categoryOrm := &model.Category{}
	err := common.SwapTo(request, categoryOrm)
	if err != nil {
		return err
	}

	err = c.CategoryDataService.UpdateCategory(categoryOrm)
	if err != nil {
		return err
	}
	response.Message = "分类更新成功"
	return nil
}

// DeleteCategory 删除品类信息
func (c *Category) DeleteCategory(ctx context.Context, request *category.DeleteCategoryRequest, response *category.DeleteCategoryResponse) error {
	err := c.CategoryDataService.DeleteCategory(request.CategoryId)
	if err != nil {
		return err
	}

	response.Message = "分类删除成功"
	return nil
}

// FindCategoryByName 根据名称查询品类信息
func (c *Category) FindCategoryByName(ctx context.Context, request *category.FindByNameRequest, response *category.CategoryResponse) error {
	categoryOrm, err := c.CategoryDataService.FindCategoryByName(request.CategoryName)
	if err != nil {
		return err
	}

	err = common.SwapTo(categoryOrm, response)
	if err != nil {
		return err
	}

	return nil
}

// FindCategoryById 根据id查询品类信息
func (c *Category) FindCategoryById(ctx context.Context, request *category.FindByIdRequest, response *category.CategoryResponse) error {
	categoryOrm, err := c.CategoryDataService.FindCategoryByID(request.CategoryId)
	if err != nil {
		return err
	}

	err = common.SwapTo(categoryOrm, response)
	if err != nil {
		return err
	}

	return nil
}

// FindCategoryByLevel 根据层级查询品类信息
func (c *Category) FindCategoryByLevel(ctx context.Context, request *category.FindByLevelRequest, response *category.FindAllResponse) error {
	categoryOrmCollection, err := c.CategoryDataService.FindCategoryByLevel(request.Level)
	if err != nil {
		return err
	}
	categoriesToResponse(categoryOrmCollection, response)
	return nil
}

// FindCategoryByParent 根据父品类查询品类信息
func (c *Category) FindCategoryByParent(ctx context.Context, request *category.FindByParentRequest, response *category.FindAllResponse) error {
	categoryOrmCollection, err := c.CategoryDataService.FindCategoryByParent(request.ParentId)
	if err != nil {
		return err
	}

	categoriesToResponse(categoryOrmCollection, response)
	return nil
}

// FindAllCategory 查询所有品类信息
func (c *Category) FindAllCategory(ctx context.Context, request *category.FindAllRequest, response *category.FindAllResponse) error {
	categoryOrmCollection, err := c.CategoryDataService.FindAllCategory()
	if err != nil {
		return err
	}

	categoriesToResponse(categoryOrmCollection, response)
	return nil
}

// categoriesToResponse 将ORM集合转换为1个FindAllResponse的响应
func categoriesToResponse(categories []model.Category, response *category.FindAllResponse) {
	for _, categoryOrm := range categories {
		categoryResponse := &category.CategoryResponse{}
		err := common.SwapTo(categoryOrm, categoryResponse)
		if err != nil {
			log.Error(err)
			break
		}
		response.Category = append(response.Category, categoryResponse)
	}
}
```