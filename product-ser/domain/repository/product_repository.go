package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"product-ser/domain/model"
)

type IProductRepository interface {
	Page(length, pageIndex int32) ([]*model.Product, error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{mysqlDB: db}
}

type ProductRepository struct {
	mysqlDB *gorm.DB
}

func (p *ProductRepository) Page(length int32, pageIndex int32) ([]*model.Product, error) {
	arr := make([]*model.Product, length)

	if length > 0 && pageIndex > 0 {
		p.mysqlDB = p.mysqlDB.Limit(int(length)).Offset((int(pageIndex) - 1) * int(length))
		if err := p.mysqlDB.Find(&arr).Error; err != nil {
			return nil, errors.Wrap(err, "query product error")
		}
		return arr, nil
	}

	return arr, errors.New("参数不匹配")
}
