package service

import (
	"product-ser/domain/model"
	"product-ser/domain/repository"
)

type IProductDataService interface {
	Page(length, pageIndex int32) (int64, []*model.Product, error)
}

func NewProductDataService(repo repository.IProductRepository) IProductDataService {
	return &ProductDataService{
		productRepo: repo,
	}
}

type ProductDataService struct {
	productRepo repository.IProductRepository
}

func (u *ProductDataService) Page(length, pageIndex int32) (int64, []*model.Product, error) {
	return u.productRepo.Page(length, pageIndex)
}
