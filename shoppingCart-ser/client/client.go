package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/consul"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"log"
	"shoppingCart-ser/common"
	"shoppingCart-ser/proto"
	"strconv"
)

var (
	client proto.AddCartService
)

func main() {
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

	rpcServer := micro.NewService(
		micro.Registry(consulReg), // 服务发现
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
	)

	client = proto.NewAddCartService("shop-cart", rpcServer.Client())

	// 具体的路由
	// 增加购物车
	router.GET("/increase", AddCart)

	service := web.NewService(
		web.Address(":6668"),
		web.Name("shop-cart-client"),
		web.Registry(consulReg),
		web.Handler(router),
	)
	service.Run()
}

func AddCart(c *gin.Context) {
	number, _ := strconv.Atoi(c.Request.FormValue("number"))
	productId, _ := strconv.Atoi(c.Request.FormValue("productId"))
	productSkuId, _ := strconv.Atoi(c.Request.FormValue("productSkuId"))

	req := &proto.AddCartReq{
		Number:       int32(number),
		ProductId:    int32(productId),
		ProductSkuId: int32(productSkuId),
	}

	resp, err := client.AddCart(context.TODO(), req)
	// 根据响应输出
	if err != nil {
		log.Println(err.Error())
		common.RespFail(c.Writer, "增加购物车失败", resp)
		return
	}

	common.RespOK(c.Writer, "增加购物车成功", resp)
}
