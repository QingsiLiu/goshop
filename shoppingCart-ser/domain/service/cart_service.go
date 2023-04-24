package service

import (
	"shoppingCart-ser/domain/model"
	"shoppingCart-ser/domain/repository"
	"shoppingCart-ser/proto"
)

type ICartService interface {
	AddCart(*proto.AddCartReq) (obj *model.ShoppingCart, err error)
	UpdateCart(req *proto.UpdateCartReq) (*model.ShoppingCart, error)
	GetOrderTotal(int32List []int32) (obj float32, err error)
	FindCart(*proto.FindCartReq) (obj *model.ShoppingCart, err error)
}

func NewCartService(repo repository.ICartRepository) ICartService {
	return &CartService{
		cartRepo: repo,
	}
}

type CartService struct {
	cartRepo repository.ICartRepository
}

func (u *CartService) GetOrderTotal(int32List []int32) (obj float32, err error) {
	return u.cartRepo.GetOrderTotal(int32List)
}

func (u *CartService) FindCart(req *proto.FindCartReq) (obj *model.ShoppingCart, err error) {
	return u.cartRepo.FindCart(req)
}

func (u *CartService) AddCart(req *proto.AddCartReq) (obj *model.ShoppingCart, err error) {
	return u.cartRepo.AddCart(req)
}

func (u *CartService) UpdateCart(req *proto.UpdateCartReq) (*model.ShoppingCart, error) {
	return u.cartRepo.UpdateCart(req)
}
