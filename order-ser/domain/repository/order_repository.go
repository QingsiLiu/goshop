package repository

import (
	"common"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"order-ser/domain/model"
	"order-ser/proto"
	"strconv"
	"strings"
	"time"
)

type IOrderRepository interface {
	FindOrder(req *proto.FindOrderReq) (*model.TraderOrder, error)
	AddTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error)
	UpdateTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error)
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{mysqlDB: db}
}

type OrderRepository struct {
	mysqlDB *gorm.DB
}

func (o *OrderRepository) UpdateTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error) {
	trade := model.TraderOrder{}
	trade.ID = req.TradeOrder.ID
	trade.OrderStatus = req.TradeOrder.OrderStatus
	trade.IsDeleted = req.TradeOrder.IsDeleted
	trade.UpdateTime = time.Now() //
	tb := o.mysqlDB.Model(&model.TraderOrder{}).
		Where("id = ?", trade.ID).
		Updates(&trade)
	fmt.Println("repository UpdateTradeOrder   >>>> ", trade)
	return &trade, tb.Error //err
}

func (o *OrderRepository) FindOrder(req *proto.FindOrderReq) (*model.TraderOrder, error) {
	id := req.GetId()
	no := req.GetOrderNo()
	obj := &model.TraderOrder{}
	tb := o.mysqlDB.Where("id = ? or order_no = ?", id, no).Find(obj)
	fmt.Println("FindTradeOrder>>>>>>> ", obj)
	return obj, tb.Error
}

func (o *OrderRepository) AddTradeOrder(req *proto.AddTradeOrderReq) (obj *model.TraderOrder, err error) {
	trade := &model.TraderOrder{}
	err = common.SwapToStruct(req, trade)
	if err != nil {
		log.Println("SwapToStruct  err :", err)
	}
	log.Println("SwapToStruct  trade :", trade)
	now := time.Now()
	trade.CreateTime = now
	trade.SubmitTime = now
	tp, _ := time.ParseDuration("30m")
	trade.ExpireTime = now.Add(tp)
	trade.OrderNo = getOrderNo(now, trade.UserId)

	tb := o.mysqlDB.Create(&trade)
	fmt.Println("repository AddTradeOrder   >>>> ", trade)
	return trade, tb.Error //err
}

// 生产 订单号   Y2022 06 27 11 00 53 948 97 103564
//
//	年    月 日 时  分 秒 毫秒 ID  随机数
func getOrderNo(time2 time.Time, userID int32) string {
	var tradeNo string
	tempNum := strconv.Itoa(rand.Intn(999999-100000+1) + 100000)
	tradeNo = "Y" + time2.Format("20060102150405.000") + strconv.Itoa(int(userID)) + tempNum
	tradeNo = strings.Replace(tradeNo, ".", "", -1)
	return tradeNo
}
