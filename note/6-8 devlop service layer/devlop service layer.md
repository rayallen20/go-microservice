# devlop service layer

## PART1. 修改获取指定用户购物车的`CartDataService.FindAllCart()`方法

注:此方法是代码生成时自带的,需要修改

### 1.1 修改方法定义

`/cart/domain/service/cart_data_service.go`:

```go
type ICartDataService interface {
	AddCart(*model.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*model.Cart) error
	FindCartByID(int64) (*model.Cart, error)
	// FindAllCart 获取指定用户的购物车
	FindAllCart(int64) ([]model.Cart, error)
}
```

### 1.2 修改方法实现

`/cart/domain/service/cart_data_service.go`:

```go
// FindAllCart 获取指定用户的购物车
func (u *CartDataService) FindAllCart(userID int64) ([]model.Cart, error) {
	return u.CartRepository.FindAll(userID)
}
```

## PART2. 实现清空指定用户购物车的`CartDataService.CleanCart()`方法

### 2.1 方法定义

`/cart/domain/service/cart_data_service.go`:

```go
type ICartDataService interface {
	AddCart(*model.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*model.Cart) error
	FindCartByID(int64) (*model.Cart, error)
	// FindAllCart 获取指定用户的购物车
	FindAllCart(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error
}
```

### 2.2 方法实现

`/cart/domain/service/cart_data_service.go`:

```go
// CleanCart 清空指定用户的购物车
func (u *CartDataService) CleanCart(userID int64) error {
	return u.CartRepository.CleanCart(userID)
}
```

## PART3. 实现向指定购物车中添加指定数量商品的`CartDataService.IncrNum()`方法

### 3.1 方法定义

`/cart/domain/service/cart_data_service.go`:

```go
type ICartDataService interface {
	AddCart(*model.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*model.Cart) error
	FindCartByID(int64) (*model.Cart, error)
	// FindAllCart 获取指定用户的购物车
	FindAllCart(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error

	// IncrNum 向指定购物车中添加指定数量的商品
	IncrNum(int64, int64) error
}
```

### 3.2 方法实现

`/cart/domain/service/cart_data_service.go`:

```go
// IncrNum 向指定购物车中添加指定数量的商品
func (u *CartDataService) IncrNum(cartID int64, num int64) error {
	return u.CartRepository.IncrNum(cartID, num)
}
```

## PART4. 实现从指定购物车中减少指定数量商品的`CartDataService.DecrNum()`方法

### 4.1 方法定义

`/cart/domain/service/cart_data_service.go`:

```go
type ICartDataService interface {
	AddCart(*model.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*model.Cart) error
	FindCartByID(int64) (*model.Cart, error)
	// FindAllCart 获取指定用户的购物车
	FindAllCart(int64) ([]model.Cart, error)

	// CleanCart 清空指定用户的购物车
	CleanCart(int64) error

	// IncrNum 向指定购物车中添加指定数量的商品
	IncrNum(int64, int64) error

	// DecrNum 从指定购物车中减少指定数量的商品
	DecrNum(int64, int64) error
}
```

### 4.2 方法实现

`/cart/domain/service/cart_data_service.go`:

```go
// DecrNum 从指定购物车中减少指定数量的商品
func (u *CartDataService) DecrNum(cartID int64, num int64) error {
	return u.CartRepository.DecrNum(cartID, num)
}
```
