package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"product-ser/domain/model"
	"product-ser/proto"
)

type IProductRepository interface {
	Page(length, pageIndex int32) (int64, []*model.Product, error)
	ShowProductDetail(int32) (obj *model.ProductDetail, err error)
	ShowProductSku(int32) (*[]model.ProductSku, error)
	ShowDetailSku(int32) (*model.ProductSku, error)
	UpdateSku(req *proto.UpdateSkuReq) (isSuccess bool, err error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{mysqlDB: db}
}

type ProductRepository struct {
	mysqlDB *gorm.DB
}

func (p *ProductRepository) Page(length int32, pageIndex int32) (int64, []*model.Product, error) {
	arr := make([]*model.Product, length)
	var count int64

	if length > 0 && pageIndex > 0 {
		p.mysqlDB = p.mysqlDB.Limit(int(length)).Offset((int(pageIndex) - 1) * int(length))
		if err := p.mysqlDB.Find(&arr).Error; err != nil {
			fmt.Println("query product err :", err)
		}
		p.mysqlDB.Model(&model.Product{}).Offset(-1).Limit(-1).Count(&count)
		return count, arr, nil
	}

	return count, arr, errors.New("参数不匹配")
}

func (u *ProductRepository) ShowProductDetail(id int32) (product *model.ProductDetail, err error) {
	sql := "select p.`id` , p.`name`, p.product_type,p.category_id ,p.starting_price, p.main_picture,\n " +
		"  pd.detail as detail ,GROUP_CONCAT(pp.picture SEPARATOR ',')  as picture_list\n" +
		"FROM `product`  p\n" +
		"  left join product_detail pd on p.id = pd.product_id\n " +
		" left join product_picture pp on p.id = pp.product_id\n " +
		" where p.`id` = ?"
	var productDetails []model.ProductDetail

	u.mysqlDB.Raw(sql, id).Scan(&productDetails)
	fmt.Println("repository ShowProductDetail   >>>> ", productDetails)
	return &productDetails[0], nil
}

func (u *ProductRepository) ShowProductSku(id int32) (product *[]model.ProductSku, err error) {
	sql := "select id ,name ,attribute_symbol_list,sell_price,stock  from product_sku where product_id= ?"
	var productSku []model.ProductSku

	u.mysqlDB.Raw(sql, id).Scan(&productSku)
	fmt.Println("repository ShowProductSku   >>>> ", productSku)
	return &productSku, nil
}

func (u *ProductRepository) ShowDetailSku(id int32) (obj *model.ProductSku, err error) {
	var productSku = &model.ProductSku{}
	u.mysqlDB.Where("id = ?", id).Find(&productSku)
	fmt.Println("repository ShowDetailSku   >>>> ", productSku)
	return productSku, nil
}

func (u *ProductRepository) UpdateSku(req *proto.UpdateSkuReq) (isSuccess bool, err error) {
	sku := req.GetProductSku()
	isSuccess = true
	//u.mysqlDB.Updates(sku)
	tb := u.mysqlDB.Debug().Model(&model.ProductSku{}).Where("id=?", sku.SkuId).Update("stock", sku.Stock)
	if tb.Error != nil {
		isSuccess = false
	}
	return isSuccess, tb.Error
}
