# devlop repository layer

## PART1. 修改创建购物车的`CartRepository.CreateCart()`方法

此处的逻辑应该是先查找,若不存在则创建.

`/cart/domain/repository/cart_repository.go`:

```go
// CreateCart 创建Cart信息
func (u *CartRepository) CreateCart(cart *model.Cart) (cartID int64, err error) {
	db := u.mysqlDb.FirstOrCreate(cart, model.Cart{
		ProductID: cart.ProductID,
		Num:       cart.Num,
		SizeID:    cart.SizeID,
		UserID:    cart.UserID,
	})

	if db.Error != nil {
		return 0, db.Error
	}

	if db.RowsAffected == 0 {
		return 0, errors.New("购物车插入失败")
	}

	return cart.ID, nil
}
```

注:该方法是生成的代码中自带的,需要修改

注:若代码报红,则在`cart/`目录下执行`go mod tidy`即可

## PART2. 修改查询指定用户购物车信息的`CartRepository.FindAll()`方法

### 2.1 修改方法定义

`/cart/domain/repository/cart_repository.go`:

```go
type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(*model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(*model.Cart) error
	
	// FindAll 获取指定用户的购物车
	FindAll(int64) ([]model.Cart, error)
}
```

### 2.2 修改方法实现

`/cart/domain/repository/cart_repository.go`:

```go
// FindAll 获取指定用户的购物车
func (u *CartRepository) FindAll(userID int64) (cartAll []model.Cart, err error) {
	return cartAll, u.mysqlDb.Where("user_id = ?", userID).Find(&cartAll).Error
}
```

注:该方法是生成的代码中自带的,需要修改

## PART3. 实现清空购物车的`CartRepository.CleanCart()`方法

### 3.1 定义方法

`/cart/domain/repository/cart_repository.go`:

```go
type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(*model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(*model.Cart) error

	// FindAll 获取指定用户的购物车
	FindAll(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error
}
```

### 3.2 实现方法

`/cart/domain/repository/cart_repository.go`:

```go
// CleanCart 清空指定用户的购物车
func (u *CartRepository) CleanCart(userID int64) error {
	return u.mysqlDb.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}
```

## PART4. 实现向指定购物车中添加指定数量商品的`CartRepository.IncrNum()`方法

### 4.1 定义方法

`/cart/domain/repository/cart_repository.go`:

```go
type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(*model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(*model.Cart) error

	// FindAll 获取指定用户的购物车
	FindAll(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error

	// IncrNum 向指定购物车中添加指定数量的商品
	IncrNum(int64, int64) error
}
```

### 4.2 实现方法

`/cart/domain/repository/cart_repository.go`:

```go
// IncrNum 向指定购物车中添加指定数量的商品
func (u *CartRepository) IncrNum(cartID int64, num int64) error {
	cartOrm := &model.Cart{ID: cartID}
	return u.mysqlDb.Model(cartOrm).UpdateColumn("num", gorm.Expr("num + ?", num)).Error
}
```

## PART5. 实现从指定购物车中减少指定数量商品的`CartRepository.DecrNum()`方法

### 5.1 定义方法

`/cart/domain/repository/cart_repository.go`:

```go
type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(*model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(*model.Cart) error

	// FindAll 获取指定用户的购物车
	FindAll(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error

	// IncrNum 向指定购物车中添加指定数量的商品
	IncrNum(int64, int64) error

	// DecrNum 从指定购物车中减少指定数量的商品
	DecrNum(int64, int64) error
}
```

### 5.2 实现方法

`/cart/domain/repository/cart_repository.go`:

```go
// DecrNum 从指定购物车中减少指定数量的商品
func (u *CartRepository) DecrNum(cartID int64, num int64) error {
	cartOrm := &model.Cart{ID: cartID}
	// 因为不能把商品减少到负数 所以要查询购物车中商品数量 >= 指定数量的购物车信息
	db := u.mysqlDb.Model(cartOrm).Where("num >= ?", num).UpdateColumn("num", gorm.Expr("num - ?", num))
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("减少失败")
	}

	return nil
}
```

## PART6. 完整代码

完整的`/cart/domain/repository/cart_repository.go`代码如下:

```go
package repository

import (
	"errors"
	"git.imooc.com/cap1573/cart/domain/model"
	"github.com/jinzhu/gorm"
)

type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(*model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(*model.Cart) error

	// FindAll 获取指定用户的购物车
	FindAll(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error

	// IncrNum 向指定购物车中添加指定数量的商品
	IncrNum(int64, int64) error

	// DecrNum 从指定购物车中减少指定数量的商品
	DecrNum(int64, int64) error
}

// NewCartRepository 创建cartRepository
func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{mysqlDb: db}
}

type CartRepository struct {
	mysqlDb *gorm.DB
}

// InitTable 初始化表
func (u *CartRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Cart{}).Error
}

// FindCartByID 根据ID查找Cart信息
func (u *CartRepository) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	cart = &model.Cart{}
	return cart, u.mysqlDb.First(cart, cartID).Error
}

// CreateCart 创建Cart信息
func (u *CartRepository) CreateCart(cart *model.Cart) (cartID int64, err error) {
	db := u.mysqlDb.FirstOrCreate(cart, model.Cart{
		ProductID: cart.ProductID,
		Num:       cart.Num,
		SizeID:    cart.SizeID,
		UserID:    cart.UserID,
	})

	if db.Error != nil {
		return 0, db.Error
	}

	if db.RowsAffected == 0 {
		return 0, errors.New("购物车插入失败")
	}

	return cart.ID, nil
}

// DeleteCartByID 根据ID删除Cart信息
func (u *CartRepository) DeleteCartByID(cartID int64) error {
	return u.mysqlDb.Where("id = ?", cartID).Delete(&model.Cart{}).Error
}

// UpdateCart 更新Cart信息
func (u *CartRepository) UpdateCart(cart *model.Cart) error {
	return u.mysqlDb.Model(cart).Update(cart).Error
}

// FindAll 获取指定用户的购物车
func (u *CartRepository) FindAll(userID int64) (cartAll []model.Cart, err error) {
	return cartAll, u.mysqlDb.Where("user_id = ?", userID).Find(&cartAll).Error
}

// CleanCart 清空指定用户的购物车
func (u *CartRepository) CleanCart(userID int64) error {
	return u.mysqlDb.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}

// IncrNum 向指定购物车中添加指定数量的商品
func (u *CartRepository) IncrNum(cartID int64, num int64) error {
	cartOrm := &model.Cart{ID: cartID}
	return u.mysqlDb.Model(cartOrm).UpdateColumn("num", gorm.Expr("num + ?", num)).Error
}

// DecrNum 从指定购物车中减少指定数量的商品
func (u *CartRepository) DecrNum(cartID int64, num int64) error {
	cartOrm := &model.Cart{ID: cartID}
	// 因为不能把商品减少到负数 所以要查询购物车中商品数量 >= 指定数量的购物车信息
	db := u.mysqlDb.Model(cartOrm).Where("num >= ?", num).UpdateColumn("num", gorm.Expr("num - ?", num))
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("减少失败")
	}

	return nil
}
```