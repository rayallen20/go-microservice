# devlop repository layer

注:编写前在`product`(项目根路径)下执行`go mod tidy`命令

## PART1. 修改`ProductRepository.InitTable()`方法

该方法用于初始化表.本案例中共有4张表,生成的代码中只创建了1张表,故需修改.

修改后的`InitTable()`方法:

```go
// InitTable 初始化表
func (u *ProductRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Product{}, &model.ProductImage{}, &model.ProductSeo{}, &model.ProductSize{}).Error
}
```

## PART2. 修改`ProductRepository.FindProductByID()`方法

由于`model.Product`结构体有外键,所以查询的时候要调用`Preload()`方法将关联的外键数据一同查询出来.

修改后的`FindProductByID()`方法:

```go
// FindProductByID 根据ID查找Product信息
func (u *ProductRepository) FindProductByID(productID int64) (product *model.Product, err error) {
	product = &model.Product{}
	return product, u.mysqlDb.Preload("ProductImage").
		Preload("ProductSeo").Preload("ProductSize").First(product, productID).Error
}
```

## PART3. 修改`ProductRepository.DeleteProductByID()`方法

由于本例中有表关联,所以删除商品需要使用事务确保删除商品时,连同商品的图片、尺码、SEO信息一同删除

修改后的`DeleteProductByID()`方法:

```go
// DeleteProductByID 根据ID删除Product信息
func (u *ProductRepository) DeleteProductByID(productID int64) error {
	// 开启事务
	tx := u.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除商品信息
	// Unscoped().Delete()表示物理删除
	err := tx.Unscoped().Where("id = ?", productID).Delete(&model.Product{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除图片信息
	err = tx.Unscoped().Where("images_product_id = ?", productID).Delete(&model.ProductImage{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除尺码信息
	err = tx.Unscoped().Where("size_product_id = ?", productID).Delete(&model.ProductSize{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除SEO信息
	err = tx.Unscoped().Where("seo_product_id = ?", productID).Delete(&model.ProductSeo{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
```

## PART4. 修改`ProductRepository.FindAll()`方法

和PART2处相同,由于`model.Product`结构体有外键,所以查询的时候要调用`Preload()`方法将关联的外键数据一同查询出来.

修改后的`FindAll()`方法:

```go
// FindAll 获取结果集
func (u *ProductRepository) FindAll() (productAll []model.Product, err error) {
	return productAll, u.mysqlDb.Preload("ProductImage").Preload("ProductSize").
		Preload("ProductSeo").Find(&productAll).Error
}
```

## PART5. 完整的repository层代码

`product/domain/repository/product_repository.go`:

```go
package repository

import (
	"git.imooc.com/cap1573/product/domain/model"
	"github.com/jinzhu/gorm"
)

type IProductRepository interface {
	InitTable() error
	FindProductByID(int64) (*model.Product, error)
	CreateProduct(*model.Product) (int64, error)
	DeleteProductByID(int64) error
	UpdateProduct(*model.Product) error
	FindAll() ([]model.Product, error)
}

// NewProductRepository 创建productRepository
func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{mysqlDb: db}
}

type ProductRepository struct {
	mysqlDb *gorm.DB
}

// InitTable 初始化表
func (u *ProductRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Product{}, &model.ProductImage{}, &model.ProductSeo{}, &model.ProductSize{}).Error
}

// FindProductByID 根据ID查找Product信息
func (u *ProductRepository) FindProductByID(productID int64) (product *model.Product, err error) {
	product = &model.Product{}
	return product, u.mysqlDb.Preload("ProductImage").
		Preload("ProductSeo").Preload("ProductSize").First(product, productID).Error
}

// CreateProduct 创建Product信息
func (u *ProductRepository) CreateProduct(product *model.Product) (int64, error) {
	return product.ID, u.mysqlDb.Create(product).Error
}

// DeleteProductByID 根据ID删除Product信息
func (u *ProductRepository) DeleteProductByID(productID int64) error {
	// 开启事务
	tx := u.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除商品信息
	// Unscoped().Delete()表示物理删除
	err := tx.Unscoped().Where("id = ?", productID).Delete(&model.Product{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除图片信息
	err = tx.Unscoped().Where("images_product_id = ?", productID).Delete(&model.ProductImage{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除尺码信息
	err = tx.Unscoped().Where("size_product_id = ?", productID).Delete(&model.ProductSize{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除SEO信息
	err = tx.Unscoped().Where("seo_product_id = ?", productID).Delete(&model.ProductSeo{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateProduct 更新Product信息
func (u *ProductRepository) UpdateProduct(product *model.Product) error {
	return u.mysqlDb.Model(product).Update(product).Error
}

// FindAll 获取结果集
func (u *ProductRepository) FindAll() (productAll []model.Product, err error) {
	return productAll, u.mysqlDb.Preload("ProductImage").Preload("ProductSize").
		Preload("ProductSeo").Find(&productAll).Error
}
```