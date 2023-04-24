package handler

import (
	"common"
	"context"
	"fmt"
	"order-ser/domain/service"
	"order-ser/proto"
)

type (
	OrderHandler struct {
		OrderService service.IOrderService
	}
)

func (h *OrderHandler) FindOrder(ctx context.Context, req *proto.FindOrderReq, resp *proto.FindOrderResp) error {
	obj, err := h.OrderService.FindOrder(req)
	if err != nil {
		fmt.Println("find order error: ", err)
	} else {
		order := &proto.TradeOrder{}
		err := common.SwapToStruct(obj, order)
		if err != nil {
			fmt.Println("转换失败", err)
		}
		resp.TradeOrder = order
	}
	return err
}

func (h *OrderHandler) AddTradeOrder(ctx context.Context, req *proto.AddTradeOrderReq, resp *proto.AddTradeOrderResp) error {
	obj, err := h.OrderService.AddTradeOrder(req)
	if err != nil {
		fmt.Println("  AddTradeOrder err :", err)
	} else {
		fmt.Println(obj.UpdateTime)
		fmt.Println(" AddTradeOrder  handler  >>>>>>  ", resp)
	}
	return err
}

func (h *OrderHandler) UpdateTradeOrder(ctx context.Context, req *proto.AddTradeOrderReq, resp *proto.AddTradeOrderResp) error {
	obj, err := h.OrderService.UpdateTradeOrder(req)
	if err != nil {
		println("  UpdateTradeOrder err :", err)
	} else {
		fmt.Println(obj.UpdateTime)
		fmt.Println(" UpdateTradeOrder  handler  >>>>>>  ", resp)
	}
	return err
}
