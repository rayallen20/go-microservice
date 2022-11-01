# interact with database

注:写代码之前在`category`目录下执行`go mod tidy`

## PART1. 定义ORM

`user/git.imooc.com/cap1573/category/domain/model/category.go`:

ORM`Category`新增字段定义:

- `CategoryName`:品类名称
- `CategoryLevel`:品类层级
- `CategoryParent`:父级品类ID
- `CategoryImage`:品类图片URL
- `CategoryDescription`:品类描述

```go
package model

type Category struct {
	// ID 主键自增ID
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	// CategoryName 品类名称
	CategoryName string `gorm:"unique_index,not_null" json:"category_name"`
	// CategoryLevel 品类层级
	CategoryLevel uint32 `json:"category_level"`
	// CategoryParent 父级品类ID
	CategoryParent int64 `json:"category_parent"`
	// CategoryImage 品类图片URL
	CategoryImage string `json:"category_image"`
	// CategoryDescription 品类描述
	CategoryDescription string `json:"category_description"`
}
```

## PART2. 定义对ORM的操作

`user/git.imooc.com/cap1573/category/domain/repository/category_repository.go`:

接口`ICategoryRepository`新增方法定义:

- `FindCategoryByName()`:根据品类名称查找品类信息
- `FindCategoryByLevel()`:根据层级查找品类信息集合
- `FindCategoryByParent()`:根据父级品类ID查找品类信息集合

接口`ICategoryRepository`的实现`CategoryRepository`新增方法实现:

- `FindCategoryByName()`
- `FindCategoryByLevel()`
- `FindCategoryByParent()`

```go
package repository

import (
	"git.imooc.com/cap1573/category/domain/model"
	"github.com/jinzhu/gorm"
)

type ICategoryRepository interface {
	InitTable() error
	FindCategoryByID(int64) (*model.Category, error)
	CreateCategory(*model.Category) (int64, error)
	DeleteCategoryByID(int64) error
	UpdateCategory(*model.Category) error
	FindAll() ([]model.Category, error)
	// FindCategoryByName 根据品类名称查找品类信息
	FindCategoryByName(string) (*model.Category, error)
	// FindCategoryByLevel 根据层级查找品类信息集合
	FindCategoryByLevel(uint32) ([]model.Category, error)
	// FindCategoryByParent 根据父级品类ID查找品类信息集合
	FindCategoryByParent(int64) ([]model.Category, error)
}

// NewCategoryRepository 创建
func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{mysqlDb: db}
}

type CategoryRepository struct {
	mysqlDb *gorm.DB
}

// InitTable 初始化表
func (u *CategoryRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Category{}).Error
}

// FindCategoryByID 根据ID查找Category信息
func (u *CategoryRepository) FindCategoryByID(categoryID int64) (category *model.Category, err error) {
	category = &model.Category{}
	return category, u.mysqlDb.First(category, categoryID).Error
}

// CreateCategory 创建Category信息
func (u *CategoryRepository) CreateCategory(category *model.Category) (int64, error) {
	return category.ID, u.mysqlDb.Create(category).Error
}

// DeleteCategoryByID 根据ID删除Category信息
func (u *CategoryRepository) DeleteCategoryByID(categoryID int64) error {
	return u.mysqlDb.Where("id = ?", categoryID).Delete(&model.Category{}).Error
}

// UpdateCategory 更新Category信息
func (u *CategoryRepository) UpdateCategory(category *model.Category) error {
	return u.mysqlDb.Model(category).Update(category).Error
}

// FindAll 获取结果集
func (u *CategoryRepository) FindAll() (categoryAll []model.Category, err error) {
	return categoryAll, u.mysqlDb.Find(&categoryAll).Error
}

// FindCategoryByName 根据层级查找品类信息集合
func (u *CategoryRepository) FindCategoryByName(categoryName string) (category *model.Category, err error) {
	category = &model.Category{}
	return category, u.mysqlDb.Where("category_name = ?", categoryName).Find(category).Error
}

// FindCategoryByLevel 根据层级查找品类信息集合
func (u *CategoryRepository) FindCategoryByLevel(level uint32) (categories []model.Category, err error) {
	return categories, u.mysqlDb.Where("category_level = ?", level).Find(categories).Error
}

// FindCategoryByParent 根据父级品类ID查找品类信息集合
func (u *CategoryRepository) FindCategoryByParent(parent int64) (categories []model.Category, err error) {
	return categories, u.mysqlDb.Where("category_parent = ?", parent).Find(categories).Error
}
```

## PART3. 定义业务逻辑

`user/git.imooc.com/cap1573/category/domain/service/category_data_service.go`:

接口`ICategoryDataService`新增方法定义:

- `FindCategoryByName()`:根据品类名称查找品类信息
- `FindCategoryByLevel()`:根据层级查找品类信息集合
- `FindCategoryByParent()`:根据父级品类ID查找品类信息集合

接口`ICategoryDataService`的实现`CategoryDataService`新增方法实现:

- `FindCategoryByName()`
- `FindCategoryByLevel()`
- `FindCategoryByParent()`

```go
package service

import (
	"git.imooc.com/cap1573/category/domain/model"
	"git.imooc.com/cap1573/category/domain/repository"
)

type ICategoryDataService interface {
	AddCategory(*model.Category) (int64, error)
	DeleteCategory(int64) error
	UpdateCategory(*model.Category) error
	FindCategoryByID(int64) (*model.Category, error)
	FindAllCategory() ([]model.Category, error)
	// FindCategoryByName 根据品类名称查找品类信息
	FindCategoryByName(string) (*model.Category, error)
	// FindCategoryByLevel 根据层级查找品类信息集合
	FindCategoryByLevel(uint32) ([]model.Category, error)
	// FindCategoryByParent 根据父级品类ID查找品类信息集合
	FindCategoryByParent(int64) ([]model.Category, error)
}

// NewCategoryDataService 创建
func NewCategoryDataService(categoryRepository repository.ICategoryRepository) ICategoryDataService {
	return &CategoryDataService{categoryRepository}
}

type CategoryDataService struct {
	CategoryRepository repository.ICategoryRepository
}

// AddCategory 插入
func (u *CategoryDataService) AddCategory(category *model.Category) (int64, error) {
	return u.CategoryRepository.CreateCategory(category)
}

// DeleteCategory 删除
func (u *CategoryDataService) DeleteCategory(categoryID int64) error {
	return u.CategoryRepository.DeleteCategoryByID(categoryID)
}

// UpdateCategory 更新
func (u *CategoryDataService) UpdateCategory(category *model.Category) error {
	return u.CategoryRepository.UpdateCategory(category)
}

// FindCategoryByID 查找
func (u *CategoryDataService) FindCategoryByID(categoryID int64) (*model.Category, error) {
	return u.CategoryRepository.FindCategoryByID(categoryID)
}

// FindAllCategory 查找
func (u *CategoryDataService) FindAllCategory() ([]model.Category, error) {
	return u.CategoryRepository.FindAll()
}

// FindCategoryByName 根据品类名称查找品类信息
func (u *CategoryDataService) FindCategoryByName(categoryName string) (category *model.Category, err error) {
	return u.CategoryRepository.FindCategoryByName(categoryName)
}

// FindCategoryByLevel 根据层级查找品类信息集合
func (u *CategoryDataService) FindCategoryByLevel(level uint32) (categories []model.Category, err error) {
	return u.CategoryRepository.FindCategoryByLevel(level)
}

// FindCategoryByParent 根据父级品类ID查找品类信息集合
func (u *CategoryDataService) FindCategoryByParent(parent int64) (categories []model.Category, err error) {
	return u.CategoryRepository.FindCategoryByParent(parent)
}
```