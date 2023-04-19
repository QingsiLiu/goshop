package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/consul"
	roundrobin "github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"log"
	"net"
	"net/http"
	"shoppingCart-ser/common"
	"shoppingCart-ser/proto"
	"strconv"
)

var (
	AddCartclient       proto.AddCartService
	ProductDetailclient proto.ShowProductDetailService
	ProductSkuclient    proto.ShowProductSkuService
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

	// 熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		http.ListenAndServe(net.JoinHostPort("101.34.10.3", "9096"), hystrixStreamHandler)
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

	AddCartclient = proto.NewAddCartService("shop-cart", rpcServer.Client())
	ProductDetailclient = proto.NewShowProductDetailService("shop-product", rpcServer.Client())
	ProductSkuclient = proto.NewShowProductSkuService("shop-product", rpcServer.Client())

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

func AddCart(c *gin.Context) {
	number, _ := strconv.Atoi(c.Request.FormValue("number"))
	productId, _ := strconv.Atoi(c.Request.FormValue("productId"))
	productSkuId, _ := strconv.Atoi(c.Request.FormValue("productSkuId"))

	req := &proto.AddCartReq{
		Number:       int32(number),
		ProductId:    int32(productId),
		ProductSkuId: int32(productSkuId),
	}
	resp := &proto.AddCartResp{}

	// 调用商品服务
	reqDetail := &proto.ProductDetailReq{Id: int32(productId)}
	respDetail, err := ProductDetailclient.ShowProductDetail(context.TODO(), reqDetail)
	if err != nil {
		log.Println(err.Error())
		common.RespFail(c.Writer, "查询商品详情失败", resp)
		return
	}
	if respDetail != nil {
		req.ProductName = respDetail.ProductDetail[0].Name
		req.ProductMainPicture = respDetail.ProductDetail[0].MainPicture
	}

	// 调用商品规格服务

	// 添加购物车，调用远程服务
	resp, err = AddCartclient.AddCart(context.TODO(), req)
	log.Println("/increase: ", resp)
	// 根据响应输出
	if err != nil {
		log.Println(err.Error())
		common.RespFail(c.Writer, "增加购物车失败", resp)
		return
	}

	common.RespOK(c.Writer, "增加购物车成功", resp)
}
