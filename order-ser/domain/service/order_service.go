package service

import (
	"order-ser/domain/model"
	"order-ser/domain/repository"
	"order-ser/proto"
)

type IOrderService interface {
	FindOrder(req *proto.FindOrderReq) (*model.TraderOrder, error)
	AddTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error)
	UpdateTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error)
}

func NewOrderService(repo repository.IOrderRepository) IOrderService {
	return &OrderService{
		orderRepo: repo,
	}
}

type OrderService struct {
	orderRepo repository.IOrderRepository
}

func (o *OrderService) FindOrder(req *proto.FindOrderReq) (*model.TraderOrder, error) {
	return o.orderRepo.FindOrder(req)
}

func (o *OrderService) AddTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error) {
	return o.orderRepo.AddTradeOrder(req)
}

func (o *OrderService) UpdateTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error) {
	return o.orderRepo.UpdateTradeOrder(req)
}
