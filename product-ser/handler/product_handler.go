package handler

import (
	"context"
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
	count, productInfo, err := p.ProductDataService.Page(request.GetLength(), request.GetPageIndex())
	if err != nil {
		return errors.Wrap(err, "page product error")
	}
	log.Println("count >>> ", count)
	log.Println(">>>>>>>>> page product success: ", productInfo)

	response.Rows = int64(request.GetLength())
	response.Total = count
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
		log.Println(">>>>>>>>>>>>> ", product)
		resp.Product = append(resp.Product, product)
	}
	return nil
}
