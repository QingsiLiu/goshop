package repository

import (
	"fmt"
	"gorm.io/gorm"
	"shoppingCart-ser/domain/model"
	"shoppingCart-ser/proto"
	"time"
)

type ICartRepository interface {
	AddCart(*proto.AddCartReq) (obj *model.ShoppingCart, err error)
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{mysqlDB: db}
}

type CartRepository struct {
	mysqlDB *gorm.DB
}

func (c *CartRepository) AddCart(req *proto.AddCartReq) (obj *model.ShoppingCart, err error) {
	cart := model.ShoppingCart{
		Number:             req.Number,
		ProductId:          req.ProductId,
		ProductSkuId:       req.ProductSkuId,
		ProductName:        req.ProductName,
		ProductMainPicture: req.ProductMainPicture,
		UserId:             req.UserId,
		CreateUser:         req.CreateUser,
	}
	cart.CreateTime = time.Now()
	tb := c.mysqlDB.Create(&cart)
	fmt.Println("repository AddCart   >>>> ", cart)
	return &cart, tb.Error
}
