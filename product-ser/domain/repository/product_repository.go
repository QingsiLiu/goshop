package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"product-ser/domain/model"
)

type IProductRepository interface {
	Page(length, pageIndex int32) (int64, []*model.Product, error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{mysqlDB: db}
}

type ProductRepository struct {
	mysqlDB *gorm.DB
}

func (p *ProductRepository) Page(length int32, pageIndex int32) (int64, []*model.Product, error) {
	arr := make([]*model.Product, length)
	var count int64

	if length > 0 && pageIndex > 0 {
		p.mysqlDB = p.mysqlDB.Limit(int(length)).Offset((int(pageIndex) - 1) * int(length))
		if err := p.mysqlDB.Find(&arr).Error; err != nil {
			fmt.Println("query product err :", err)
		}
		p.mysqlDB.Model(&model.Product{}).Offset(-1).Limit(-1).Count(&count)
		return count, arr, nil
	}

	return count, arr, errors.New("参数不匹配")
}
