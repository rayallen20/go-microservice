# devlop model layer

## PART1. 商品图片ORM定义

在`domain/model/`下新建`product_image.go`:

```go
package model

// ProductImage 商品图片信息ORM
type ProductImage struct {
	// ID 商品图片id
	ID int64 `gorm:"primary_key;not_null:auto_increment" json:"id"`
	
	// ImageName 图片名称
	ImageName string `json:"image_name"`
	
	// ImageCode 图片编码 用于确保添加/删除操作的幂等性
	// 因为分布式系统中 有可能受网络抖动等原因导致重复发送请求
	// 故需要有一个唯一标识来保证同样的请求不会重复进行操作
	ImageCode string `gorm:"unique_index;not_null" json:"image_code"`
	
	// ImageUrl 图片url
	ImageUrl string `json:"image_url"`
	
	// ImageProductID 外键 图片所属商品id
	ImageProductID int64 `json:"image_product_id"`
}
```

## PART2. 商品尺码ORM定义

在`domain/model/`下新建`product_size.go`:

```go
package model

// ProductSize 商品尺码信息ORM
type ProductSize struct {
	// ID 商品尺码id
	ID int64 `gorm:"primary_key;not_null:auto_increment" json:"id"`

	// SizeName 商品尺码名称
	SizeName string `json:"size_name"`

	// SizeCode 商品尺码编码 用于确保添加/删除操作的幂等性
	SizeCode string `gorm:"unique_index;not_null" json:"size_code"`

	// SizeProductID 外键 商品尺码所属商品id
	SizeProductID int64 `json:"size_product_id"`
}
```

## PART3. 商品SEO信息ORM定义

在`domain/model/`下新建`product_seo.go`:

```go
package model

// ProductSeo 商品搜索引擎优化信息ORM
type ProductSeo struct {
	// ID 商品seo信息id
	ID int64 `gorm:"primary_key;not_null:auto_increment" json:"id"`

	// SeoTitle 商品seo信息标题
	SeoTitle string `json:"seo-title"`

	// SeoKeywords 商品seo信息关键字
	SeoKeywords string `json:"seo-keywords"`

	// SeoCode 商品seo编码 用于确保添加/删除操作的幂等性
	SeoCode string `gorm:"unique_index;not_null" json:"seo_code"`

	// SeoProductID 外键 商品seo信息所属商品id
	SeoProductID int64 `json:"seo_product_id"`
}
```

## PART4. 商品ORM定义

在`domain/model/`下编辑`product.go`:

```go
package model

// Product 商品信息ORM
type Product struct {
	// ID 商品信息ID
	ID int64 `gorm:"primary_key;not_null:auto_increment" json:"id"`

	// ProductName 商品名称
	ProductName string `json:"productName"`

	// ProductSku 商品唯一编码 类似商品的条形码
	ProductSku string `gorm:"unique_index;not_null" json:"product_sku"`

	// ProductDescription 商品描述信息
	ProductDescription string `json:"product_description"`

	// ProductImage 商品对应图片信息ORM集合
	ProductImage []ProductImage `gorm:"ForeignKey:ImageProductID" json:"product_image"`

	// ProductSize 商品对应尺码信息ORM集合
	ProductSize []ProductSize `gorm:"ForeignKey:SizeProductID" json:"product_size"`

	// ProductSeo 商品对应seo信息ORM
	ProductSeo ProductSeo `gorm:"SeoProductID" json:"product_seo"`
}
```
