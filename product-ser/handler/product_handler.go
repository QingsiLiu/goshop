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

type (
	ProductHandler struct {
		ProductDataService service.IProductDataService
	}
)

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

func (p *ProductHandler) ShowProductDetail(ctx context.Context, request *proto.ProductDetailReq, response *proto.ProductDetailResp) error {
	obj, err := p.ProductDataService.ShowProductDetail(request.GetId())
	if err != nil {
		println("ShowProductDetail   err :", err)
	}
	productDetail := &proto.ProductDetail{}
	err1 := common.SwapToStruct(obj, productDetail)
	if err1 != nil {
		println("ShowProductDetail SwapToStruct  err :", err1)
	}
	response.ProductDetail = append(response.ProductDetail, productDetail)
	return nil
}

// ShowProductSku 商品SKU列表
func (p *ProductHandler) ShowProductSku(ctx context.Context, req *proto.ProductSkuReq, resp *proto.ProductSkuResp) error {
	//count := u.ProductDataService.CountNum()
	obj, err := p.ProductDataService.ShowProductSku(req.GetProductId())
	if err != nil {
		println("ShowProductSku   err :", err)
	}
	err1 := ObjSkuForResp(obj, resp)
	if err1 != nil {
		println("ShowProductSku SwapToStruct  err :", err1)
	}
	return nil
}

func ObjSkuForResp(obj *[]model.ProductSku, resp *proto.ProductSkuResp) (err error) {
	for _, v := range *obj {
		productSku := &proto.ProductSku{}
		err := common.SwapToStruct(v, productSku)
		if err != nil {
			return err
		}
		resp.ProductSku = append(resp.ProductSku, productSku)
	}
	return nil
}

// 商品SKU详情
func (p *ProductHandler) ShowDetailSku(ctx context.Context, req *proto.ProductDetailReq, resp *proto.ProductSkuResp) error {
	//count := u.ProductDataService.CountNum()
	obj, err := p.ProductDataService.ShowDetailSku(req.Id)
	if err != nil {
		println("ShowDetailSku   err :", err)
	}
	productSku := &proto.ProductSku{}
	err = common.SwapToStruct(obj, productSku)
	if err != nil {
		return err
	}
	resp.ProductSku = append(resp.ProductSku, productSku)
	return nil
}

// 修改商品SKU
func (p *ProductHandler) UpdateSku(ctx context.Context, req *proto.UpdateSkuReq, resp *proto.UpdateSkuResp) error {
	//count := u.ProductDataService.CountNum()
	isSuccess, err := p.ProductDataService.UpdateSku(req)
	if err != nil {
		resp.IsSuccess = isSuccess
		println("UpdateSku   err :", err)
	}
	resp.IsSuccess = isSuccess
	return err
}
