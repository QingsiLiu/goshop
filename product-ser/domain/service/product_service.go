package service

import (
	"product-ser/domain/model"
	"product-ser/domain/repository"
	"product-ser/proto"
)

type IProductDataService interface {
	Page(length, pageIndex int32) (int64, []*model.Product, error)
	ShowProductDetail(int32) (obj *model.ProductDetail, err error)
	ShowProductSku(int32) (obj *[]model.ProductSku, err error)
	ShowDetailSku(int32) (obj *model.ProductSku, err error)
	UpdateSku(req *proto.UpdateSkuReq) (isSuccess bool, err error)
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

func (u *ProductDataService) ShowProductDetail(id int32) (product *model.ProductDetail, err error) {
	return u.productRepo.ShowProductDetail(id)
}

func (u *ProductDataService) ShowProductSku(id int32) (product *[]model.ProductSku, err error) {

	return u.productRepo.ShowProductSku(id)
}
func (u *ProductDataService) ShowDetailSku(id int32) (product *model.ProductSku, err error) {

	return u.productRepo.ShowDetailSku(id)
}

func (u *ProductDataService) UpdateSku(req *proto.UpdateSkuReq) (isSuccess bool, err error) {

	return u.productRepo.UpdateSku(req)
}
