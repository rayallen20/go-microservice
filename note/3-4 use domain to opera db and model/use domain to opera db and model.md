# use domain to opera db and model

## PART1. 创建模型和数据库操作接口

- step1. 在`domain/model`下新建`user.go`

该文件用于定义ORM

`user.go`内容:

```go
package model

type User struct {
	// 主键
	Id int64 `gorm:"primary_key;not_null;auto_increment"`

	// 账号
	UserName string `gorm:"unique_index;not_null"`

	// 用户名
	FirstName string

	// 密码
	HashPassword string
}
```

- step2. 在`domain/repository`下新建`user_repository.go`

该文件用于定义数据库操作接口

`user_repository.go`内容

```go
package repository

import (
	"git.imooc.com/rayallen20c/user/domain/model"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	// InitTable 初始化数据表
	InitTable() error

	// FindUserByName 根据用户名称查找信息
	FindUserByName(string) (*model.User, error)

	// FindUserById 根据用户ID查找信息
	FindUserById(int64) (*model.User, error)

	// CreateUser 创建用户
	CreateUser(*model.User) (int64, error)

	// DeleteUserById 根据用户信息id用户
	DeleteUserById(int64) error

	// UpdateUser 更新用户信息
	UpdateUser(*model.User) error

	// FindAll 查找所有用户
	FindAll() ([]*model.User, error)
}
```

## PART2. 实现数据库操作接口

在`user_repository.go`中新建结构体`UserRepository`,并实现数据库操作接口`IUserRepository`

- step1. 实现接口`IUserRepository`

```go
package repository

import (
	"git.imooc.com/rayallen20c/user/domain/model"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	// InitTable 初始化数据表
	InitTable() error

	// FindUserByName 根据用户名称查找信息
	FindUserByName(string) (*model.User, error)

	// FindUserById 根据用户ID查找信息
	FindUserById(int64) (*model.User, error)

	// CreateUser 创建用户
	CreateUser(*model.User) (int64, error)

	// DeleteUserById 根据用户信息id用户
	DeleteUserById(int64) error

	// UpdateUser 更新用户信息
	UpdateUser(*model.User) error

	// FindAll 查找所有用户
	FindAll() ([]*model.User, error)
}

// UserRepository 实现接口IUserRepository 用于操作DB
type UserRepository struct {
	mysqldb *gorm.DB
}

// InitTable 初始化表
func (u *UserRepository) InitTable() error {
	return u.mysqldb.CreateTable(&model.User{}).Error
}

// FindUserByName 根据用户名称查找信息
func (u *UserRepository) FindUserByName(name string) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqldb.Where("user_name = ?", name).Find(user).Error
}

// FindUserById 根据用户ID查找信息
func (u *UserRepository) FindUserById(id int64) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqldb.First(user, id).Error
}

// CreateUser 创建用户
func (u *UserRepository) CreateUser(user *model.User) (userId int64, err error) {
	return user.Id, u.mysqldb.Create(user).Error
}

// DeleteUserById 根据用户id删除用户
func (u *UserRepository) DeleteUserById(id int64) error {
	return u.mysqldb.Where("id = ?", id).Delete(&model.User{}).Error
}

// UpdateUser 更新用户信息
func (u *UserRepository) UpdateUser(user *model.User) error {
	return u.mysqldb.Model(user).Update(&user).Error
}

// FindAll 查找所有用户
func (u *UserRepository) FindAll() (userAll []*model.User, err error) {
	return userAll, u.mysqldb.Find(&userAll).Error
}
```

- step2. 实现一个函数,该函数用于创建接口`IUserRepository`的实例`UserRepository`

`user_repository.go`:

```go
package repository

import (
	"git.imooc.com/rayallen20c/user/domain/model"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	// InitTable 初始化数据表
	InitTable() error

	// FindUserByName 根据用户名称查找信息
	FindUserByName(string) (*model.User, error)

	// FindUserById 根据用户ID查找信息
	FindUserById(int64) (*model.User, error)

	// CreateUser 创建用户
	CreateUser(*model.User) (int64, error)

	// DeleteUserById 根据用户信息id用户
	DeleteUserById(int64) error

	// UpdateUser 更新用户信息
	UpdateUser(*model.User) error

	// FindAll 查找所有用户
	FindAll() ([]*model.User, error)
}

// UserRepository 实现接口IUserRepository 用于操作DB
type UserRepository struct {
	mysqldb *gorm.DB
}

// InitTable 初始化表
func (u *UserRepository) InitTable() error {
	return u.mysqldb.CreateTable(&model.User{}).Error
}

// FindUserByName 根据用户名称查找信息
func (u *UserRepository) FindUserByName(name string) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqldb.Where("user_name = ?", name).Find(user).Error
}

// FindUserById 根据用户ID查找信息
func (u *UserRepository) FindUserById(id int64) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqldb.First(user, id).Error
}

// CreateUser 创建用户
func (u *UserRepository) CreateUser(user *model.User) (userId int64, err error) {
	return user.Id, u.mysqldb.Create(user).Error
}

// DeleteUserById 根据用户id删除用户
func (u *UserRepository) DeleteUserById(id int64) error {
	return u.mysqldb.Where("id = ?", id).Delete(&model.User{}).Error
}

// UpdateUser 更新用户信息
func (u *UserRepository) UpdateUser(user *model.User) error {
	return u.mysqldb.Model(user).Update(&user).Error
}

// FindAll 查找所有用户
func (u *UserRepository) FindAll() (userAll []*model.User, err error) {
	return userAll, u.mysqldb.Find(&userAll).Error
}

// NewIUserRepository 创建接口IUserRepository的实例
func NewIUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{mysqldb: db}
}
```

## PART3. 定义业务逻辑接口

在`domain/service`目录下创建`user_data_service.go`

该文件用于定义业务逻辑接口`IUserDataService`

- step1. 定义接口

`user_data_service.go`:

```go
package service

import "git.imooc.com/rayallen20c/user/domain/model"

type IUserDataService interface {
	// AddUser 新增用户
	AddUser(user *model.User) (int64, error)

	// DeleteUser 删除用户信息
	DeleteUser(int64) error

	// UpdateUser 更新用户信息
	// isChangePwd: 是否更新密码
	UpdateUser(user *model.User, isChangePwd bool) (err error)

	// FindUserByName 根据账号查找用户
	FindUserByName(string) (*model.User, error)

	// CheckPwd 确认用户的账号和密码是否正确
	CheckPwd(userName, pwd string) (isOk bool, err error)
}
```

- step2. 实现接口`IUserDataService`

在`user_data_service.go`中新建结构体`UserDataService`,并实现业务逻辑接口`IUserDataService`

```go
package service

import (
	"errors"
	"git.imooc.com/rayallen20c/user/domain/model"
	"git.imooc.com/rayallen20c/user/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserDataService interface {
	// AddUser 新增用户
	AddUser(user *model.User) (int64, error)

	// DeleteUser 删除用户信息
	DeleteUser(int64) error

	// UpdateUser 更新用户信息
	// isChangePwd: 是否更新密码
	UpdateUser(user *model.User, isChangePwd bool) (err error)

	// FindUserByName 根据账号查找用户
	FindUserByName(string) (*model.User, error)

	// CheckPwd 确认用户的账号和密码是否正确
	CheckPwd(userName, pwd string) (isOk bool, err error)
}

type UserDataService struct {
	UserRepository repository.IUserRepository
}

// AddUser 新增用户
func (u *UserDataService) AddUser(user *model.User) (userId int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return user.Id, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.CreateUser(user)
}

// DeleteUser 删除用户信息
func (u *UserDataService) DeleteUser(userId int64) error {
	return u.UserRepository.DeleteUserById(userId)
}

// UpdateUser 更新用户信息
// isChangePwd: 是否更新密码
func (u *UserDataService) UpdateUser(user *model.User, isChangePwd bool) (err error) {
	if isChangePwd {
		pwdByte, genPwdErr := GeneratePassword(user.HashPassword)
		if genPwdErr != nil {
			return genPwdErr
		}
		user.HashPassword = string(pwdByte)
	}
	return u.UserRepository.UpdateUser(user)
}

// FindUserByName 根据账号查找用户
func (u *UserDataService) FindUserByName(userName string) (user *model.User, err error) {
	return u.UserRepository.FindUserByName(userName)
}

// CheckPwd 确认用户的账号和密码是否正确
func (u *UserDataService) CheckPwd(userName, pwd string) (isOk bool, err error) {
	user, err := u.UserRepository.FindUserByName(userName)
	if err != nil {
		return false, err
	}

	return ValidatePassword(pwd, user.HashPassword)
}

// GeneratePassword 加密用户密码
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 验证用户密码
func ValidatePassword(userPassword, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}
```

注:`GeneratePassword()`和`ValidatePassword()`属于辅助功能的函数.此处我觉得这2个函数应该是私有的.

- step3. 实现一个函数,该函数用于创建接口`IUserDataService`的实例`UserDataService`

```go
package service

import (
	"errors"
	"git.imooc.com/rayallen20c/user/domain/model"
	"git.imooc.com/rayallen20c/user/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserDataService interface {
	// AddUser 新增用户
	AddUser(user *model.User) (int64, error)

	// DeleteUser 删除用户信息
	DeleteUser(int64) error

	// UpdateUser 更新用户信息
	// isChangePwd: 是否更新密码
	UpdateUser(user *model.User, isChangePwd bool) (err error)

	// FindUserByName 根据账号查找用户
	FindUserByName(string) (*model.User, error)

	// CheckPwd 确认用户的账号和密码是否正确
	CheckPwd(userName, pwd string) (isOk bool, err error)
}

type UserDataService struct {
	UserRepository repository.IUserRepository
}

// AddUser 新增用户
func (u *UserDataService) AddUser(user *model.User) (userId int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return user.Id, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.CreateUser(user)
}

// DeleteUser 删除用户信息
func (u *UserDataService) DeleteUser(userId int64) error {
	return u.UserRepository.DeleteUserById(userId)
}

// UpdateUser 更新用户信息
// isChangePwd: 是否更新密码
func (u *UserDataService) UpdateUser(user *model.User, isChangePwd bool) (err error) {
	if isChangePwd {
		pwdByte, genPwdErr := GeneratePassword(user.HashPassword)
		if genPwdErr != nil {
			return genPwdErr
		}
		user.HashPassword = string(pwdByte)
	}
	return u.UserRepository.UpdateUser(user)
}

// FindUserByName 根据账号查找用户
func (u *UserDataService) FindUserByName(userName string) (user *model.User, err error) {
	return u.UserRepository.FindUserByName(userName)
}

// CheckPwd 确认用户的账号和密码是否正确
func (u *UserDataService) CheckPwd(userName, pwd string) (isOk bool, err error) {
	user, err := u.UserRepository.FindUserByName(userName)
	if err != nil {
		return false, err
	}

	return ValidatePassword(pwd, user.HashPassword)
}

// GeneratePassword 加密用户密码
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 验证用户密码
func ValidatePassword(userPassword, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}

// NewUserDataService 创建接口IUserDataService的实例
func NewUserDataService(userRepository repository.IUserRepository) IUserDataService {
	return &UserDataService{UserRepository: userRepository}
}
```

