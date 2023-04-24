package main

import (
	"common"
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/consul"
	roundrobin "github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/lithammer/shortuuid/v3"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"log"
	"net"
	"net/http"
	"order-ser/proto"
	"strconv"
)

var (
	UpdateCartClient    proto.UpdateCartService
	GetUserTokenClient  proto.GetUserTokenService
	GetOrderTotalClient proto.GetOrderTotalService
	AddOrderClient      proto.AddTradeOrderService
	UpdateOrderClient   proto.UpdateTradeOrderService
	FindCartClient      proto.FindCartService
	FindOrderClient     proto.FindOrderService
)

func main() {
	const (
		DtmServer = "http://101.34.10.3:36789/api/dtmsvr"
		QSBusi    = "http://localhost:6668"
	)

	/* var (
	    CartId int32 = 1
	    Number int32 = 1
	)*/
	resp := &proto.AddTradeOrderResp{}

	router := gin.Default()

	// 注册到consul
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"101.34.10.3:8500"}
	})

	// 初始化链路追踪的jaeger
	trancer, io, err := common.NewTraner("shop-cart-client", "101.34.10.3:6831")
	if err != nil {
		log.Fatalln(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(trancer)

	// 熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		http.ListenAndServe(net.JoinHostPort("101.34.10.3", "9097"), hystrixStreamHandler)
	}()

	rpcServer := micro.NewService(
		micro.Registry(consulReg), // 服务发现
		// 链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 加入熔断器
		micro.WrapClient(NewClientHystrixWrapper()),
		// 负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	UpdateCartClient = proto.NewUpdateCartService("shop-cart", rpcServer.Client())
	GetUserTokenClient = proto.NewGetUserTokenService("shop-user", rpcServer.Client())
	GetOrderTotalClient = proto.NewGetOrderTotalService("shop-cart", rpcServer.Client())
	AddOrderClient = proto.NewAddTradeOrderService("trade-order", rpcServer.Client())
	UpdateOrderClient = proto.NewUpdateTradeOrderService("trade-order", rpcServer.Client())
	FindCartClient = proto.NewFindCartService("shop-cart", rpcServer.Client())
	FindOrderClient = proto.NewFindOrderService("trade-order", rpcServer.Client())

	router.POST("/updateCart", func(c *gin.Context) {
		req := &proto.UpdateCartReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		req.IsDeleted = true
		_, err := UpdateCartClient.UpdateCart(context.TODO(), req)
		if err != nil {
			log.Println("/updateCart err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "删除购物车失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"updateCart": "SUCCESS", "Message": "删除购物车成功!"})
	})

	router.POST("/updateCart-compensate", func(c *gin.Context) {
		req := &proto.UpdateCartReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		req.IsDeleted = false
		_, err := UpdateCartClient.UpdateCart(context.TODO(), req)
		if err != nil {
			log.Println("/updateCart err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "回滚购物车失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"updateCart": "SUCCESS", "Message": "回滚购物车成功!"})
	})

	router.POST("/addTrade", func(c *gin.Context) {
		req := &proto.AddTradeOrderReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		_, err := AddOrderClient.AddTradeOrder(context.TODO(), req)
		if err != nil {
			log.Println("/addTrade err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "新增订单失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"addTrade": "SUCCESS", "Message": "新增订单成功!"})
	})

	router.POST("/addTrade-compensate", func(c *gin.Context) {
		req := &proto.AddTradeOrderReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		req.TradeOrder.IsDeleted = true
		_, err := AddOrderClient.AddTradeOrder(context.TODO(), req)
		if err != nil {
			log.Println("/addTrade err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "回滚订单失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"addTrade": "SUCCESS", "Message": "回滚订单成功!"})
	})

	// 具体的路由
	router.GET("/cartAdvanceOrder", func(c *gin.Context) {
		// 开始校验登录
		uuid := c.Request.Header["Uuid"][0]
		cc := common.GetInput(uuid)
		out := common.SQ(cc)
		sum := 0
		for o := range out {
			sum += o
		}
		//Token校验
		//拼接请求信息
		tokenReq := &proto.TokenReq{
			Uuid: uuid,
		}
		//响应
		tokenResp, err := GetUserTokenClient.GetUserToken(context.TODO(), tokenReq)
		if err != nil || tokenResp.IsLogin == false {
			log.Println("GetUserToken  err : ", err)
			common.RespFail(c.Writer, "未登录！", tokenResp)
			return
		}
		//拼接请求信息
		log.Println("GetUserToken success : ", tokenResp)
		//结束检验登录

		tmpStr := c.Request.FormValue("cartIds")
		cartIds := common.SplitToInt32List(tmpStr, ",")
		isVirtual, _ := strconv.ParseBool(c.Request.FormValue("isVirtual"))
		recipientAddressId, _ := strconv.Atoi(c.Request.FormValue("recipientAddressId"))

		// 统计价格
		totalReq := &proto.OrderTotalReq{CartIds: cartIds}

		// 开始校验cart
		findCartReq := &proto.FindCartReq{Id: cartIds[0]}
		cart, err := FindCartClient.FindCart(context.TODO(), findCartReq)
		if err != nil {
			log.Println("FindCart  err : ", err)
			common.RespFail(c.Writer, "查询购物车失败！", tokenResp)
			return
		}
		if cart.ShoppingCart.IsDeleted {
			common.RespFail(c.Writer, " 购物车已失效！", tokenResp)
			return
		}
		// 结束cart的状态校验

		totalPriceResp, _ := GetOrderTotalClient.GetOrderTotal(context.TODO(), totalReq)
		log.Println("totalPrice>>>>>>>>>>>>>>>>>>>>>>>>   ", totalPriceResp)

		// 构建 order
		tradeOrder := &proto.TradeOrder{}
		tradeOrder.UserId = int32(sum)
		tradeOrder.CreateUser = int32(sum)
		tradeOrder.OrderStatus = 1
		tradeOrder.TotalAmount = totalPriceResp.TotalPrice
		req := &proto.AddTradeOrderReq{
			CartIds:            cartIds,
			IsVirtual:          isVirtual,
			RecipientAddressId: int32(recipientAddressId),
			TradeOrder:         tradeOrder,
		}
		updateCartReq := &proto.UpdateCartReq{
			Id: cartIds[0],
		}

		// 全局事务
		gid := shortuuid.New()
		saga := dtmcli.NewSaga(DtmServer, gid).
			Add(QSBusi+"/updateCart", QSBusi+"/updateCart-compensate", updateCartReq).
			Add(QSBusi+"/addTrade", QSBusi+"/addTrade-compensate", req)
		err = saga.Submit()
		if err != nil {
			log.Println("sage submit err: ", err)
			common.RespFail(c.Writer, "添加失败", resp)
		}
		log.Println("/saga submit submit: ", gid)

		common.RespOK(c.Writer, "请求成功", resp)
	})

	router.POST("findOrder", func(c *gin.Context) {
		req := &proto.FindOrderReq{}
		req.Id = c.PostForm("id")
		req.OrderNo = c.PostForm("orderNo")

		obj, err := FindOrderClient.FindOrder(context.TODO(), req)
		if err != nil {
			log.Println("findOrder err :", err)
			common.RespFail(c.Writer, "查询失败", resp)
		}
		fmt.Println("findOrder:", obj)
		c.JSON(http.StatusOK, gin.H{"findOrder": "SUCCESS", "Message": obj})
	})

	service := web.NewService(
		web.Address(":6669"),
		web.Name("shop-cart-client"),
		web.Registry(consulReg),
		web.Handler(router),
	)
	service.Run()
}

type ClientWrapper struct {
	client.Client
}

func (c *ClientWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 正常执行逻辑
		fmt.Println("call success ", req.Service()+"."+req.Endpoint())
		return c.Client.Call(ctx, req, resp, opts...)
	}, func(err error) error {
		// 有异常
		fmt.Println("call error: ", err)
		return err
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &ClientWrapper{i}
	}
}
