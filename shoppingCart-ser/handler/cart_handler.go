package handler

import (
	"context"
	"github.com/pkg/errors"
	"shoppingCart-ser/common"
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
