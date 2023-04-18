package service

import (
	"shoppingCart-ser/domain/model"
	"shoppingCart-ser/domain/repository"
	"shoppingCart-ser/proto"
)

type ICartService interface {
	AddCart(*proto.AddCartReq) (obj *model.ShoppingCart, err error)
}

func NewCartService(repo repository.ICartRepository) ICartService {
	return &CartService{
		cartRepo: repo,
	}
}

type CartService struct {
	cartRepo repository.ICartRepository
}

func (u *CartService) AddCart(req *proto.AddCartReq) (obj *model.ShoppingCart, err error) {
	return u.cartRepo.AddCart(req)
}
