package model

type ProductSku struct {
	//gorm.Model
	SkuId               int32   `gorm:"column:id" json:"skuId"`
	Name                string  `json:"name"`
	AttributeSymbolList string  `gorm:"column:attribute_symbol_list" json:"attributeSymbolList"`
	SellPrice           float32 `gorm:"column:sell_price" json:"sellPrice"`
	Stock               int32   `gorm:"default:1" json:"stock"`
}

func (table *ProductSku) TableName() string {
	return "product_sku"
}
