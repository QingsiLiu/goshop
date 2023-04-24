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
	UpdateCart(req *proto.UpdateCartReq) (*model.ShoppingCart, error)
	GetOrderTotal(int32List []int32) (obj float32, err error)
	FindCart(*proto.FindCartReq) (obj *model.ShoppingCart, err error)
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

func (u *CartRepository) UpdateCart(req *proto.UpdateCartReq) (obj *model.ShoppingCart, err error) {
	cart := model.ShoppingCart{
		Number:             req.Number,
		ProductId:          req.ProductId,
		ProductSkuId:       req.ProductSkuId,
		ProductName:        req.ProductName,
		ProductMainPicture: req.ProductMainPicture,
		UserId:             req.UserId,
		ID:                 req.Id,
		IsDeleted:          req.IsDeleted,
		UpdateUser:         req.UpdateUser,
	}
	cart.UpdateTime = time.Now() //  &cart
	tb := u.mysqlDB.Model(&model.ShoppingCart{}).
		Where("id = ?", cart.ID).
		Updates(&cart)
	//Update("is_deleted", 1)
	//err = errors.New("400", "测试异常", 400)
	fmt.Println("repository UpdateCart   >>>> ", cart)
	return &cart, tb.Error //err
}

// 统计订单价格
func (u *CartRepository) GetOrderTotal(int32List []int32) (obj float32, err error) {
	sql := "select sum(c.Number * s.sell_price) from `shopping_cart` c \n" +
		"LEFT JOIN `product_sku` s on c.product_sku_id=s.id\n" +
		"where c.id in ?"
	var totalPrice float32
	tb := u.mysqlDB.Raw(sql, int32List).Scan(&totalPrice)
	//Update("order_status", trade.OrderStatus)
	//err = errors.New("400", "测试异常", 400)
	fmt.Println("GetOrderTotal   >>>> ", totalPrice)
	return totalPrice, tb.Error //err
}

// 查询购物车
func (u *CartRepository) FindCart(req *proto.FindCartReq) (obj *model.ShoppingCart, err error) {
	id := req.Id
	cart := &model.ShoppingCart{}

	tb := u.mysqlDB.Where("id = ?", id).Find(cart)
	fmt.Println("FindCart     >>>> ", cart)
	return cart, tb.Error //err
}
