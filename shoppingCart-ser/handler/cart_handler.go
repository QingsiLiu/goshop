package handler

import (
	"common"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"shoppingCart-ser/domain/service"
	"shoppingCart-ser/proto"
)

type (
	CartHandler struct {
		CartService service.ICartService
	}
)

func (h *CartHandler) AddCart(ctx context.Context, req *proto.AddCartReq, resp *proto.AddCartResp) error {
	obj, err := h.CartService.AddCart(req)
	if err != nil {
		return errors.Wrap(err, "add cart error")
	}

	cart := &proto.ShoppingCart{}
	err = common.SwapToStruct(obj, cart)
	if err != nil {
		return err
	}
	return nil
}

func (u *CartHandler) UpdateCart(ctx context.Context, req *proto.UpdateCartReq, resp *proto.UpdateCartResp) error {
	obj, err := u.CartService.UpdateCart(req)
	if err != nil {
		println("  UpdateCart err :", err)
	} else {
		resp.CanSetShoppingCartNumber = int64(obj.Number)
		resp.ShoppingCartNumber = int64(obj.Number)
		resp.IsBeyondMaxLimit = false // 查询sku
		resp.ID = obj.ID              //增加新增cart的ID
		fmt.Println(" UpdateCart  handler  >>>>>>  ", resp)
	}
	return err
}

func (u *CartHandler) FindCart(ctx context.Context, req *proto.FindCartReq, obj *proto.FindCartResp) error {
	//int32List := SplitToInt32List(req.GetCartIds(), ",")
	cart, err := u.CartService.FindCart(req)
	obj.ShoppingCart = &proto.ShoppingCart{}
	obj.ShoppingCart.Id = cart.ID
	obj.ShoppingCart.UserId = cart.UserId
	obj.ShoppingCart.IsDeleted = cart.IsDeleted
	//其他需要再加
	return err
}

// (int32List []int32) (obj *float32, err error)
func (u *CartHandler) GetOrderTotal(ctx context.Context, req *proto.OrderTotalReq, obj *proto.OrderTotalResp) error {
	obj.TotalPrice, _ = u.CartService.GetOrderTotal(req.GetCartIds())
	return nil
}
