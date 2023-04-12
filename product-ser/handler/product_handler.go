package handler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"product-ser/common"
	"product-ser/domain/model"
	"product-ser/domain/service"
	"product-ser/proto"
)

type ProductHandler struct {
	ProductDataService service.IProductDataService
}

func (p *ProductHandler) Page(ctx context.Context, request *proto.PageReq, response *proto.PageResp) error {
	productInfo, err := p.ProductDataService.Page(request.GetLength(), request.GetPageIndex())
	if err != nil {
		return errors.Wrap(err, "page product error")
	}
	log.Println(">>>>>>>>> page product success: ", productInfo)

	err = productForResp(productInfo, response)
	if err != nil {
		return err
	}

	return nil
}

func productForResp(products []*model.Product, resp *proto.PageResp) error {
	for _, v := range products {
		product := &proto.Product{}
		err := common.SwapToStruct(v, product)
		if err != nil {
			return err
		}
		fmt.Println(">>>>>>>>>>>>> ", product)
		resp.Product = append(resp.Product, product)
	}
	return nil
}
