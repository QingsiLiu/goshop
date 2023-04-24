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
	"shoppingCart-ser/proto"
	"strconv"
)

var (
	AddCartClient       proto.AddCartService
	ProductDetailclient proto.ShowProductDetailService
	ShowDetailSkuclient proto.ShowDetailSkuService
	UpdateSkuClient     proto.UpdateSkuService
	GetUserTokenClient  proto.GetUserTokenService
	UpdateCartClient    proto.UpdateCartService
)

func main() {
	const (
		DtmServer = "http://101.34.10.3:36789/api/dtmsvr"
		QSBusi    = "http://localhost:6668"
	)

	var (
		CartId int32 = 1
		Number int32 = 1
	)
	resp := &proto.AddCartResp{}

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

	AddCartClient = proto.NewAddCartService("shop-cart", rpcServer.Client())
	ProductDetailclient = proto.NewShowProductDetailService("shop-product", rpcServer.Client())
	ShowDetailSkuclient = proto.NewShowDetailSkuService("shop-product", rpcServer.Client())
	UpdateSkuClient = proto.NewUpdateSkuService("shop-product", rpcServer.Client())
	GetUserTokenClient = proto.NewGetUserTokenService("shop-user", rpcServer.Client())
	UpdateCartClient = proto.NewUpdateCartService("shop-cart", rpcServer.Client())

	// 具体的路由
	// 增加购物车
	router.GET("/increase", AddCart)

	// 开始拆分 DTM服务
	router.POST("updateSku", func(c *gin.Context) {
		req := &proto.UpdateSkuReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		_, err := UpdateSkuClient.UpdateSku(context.TODO(), req)
		if err != nil {
			log.Println("/updateSku err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "修改库存失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"updateSku": "SUCCESS", "Message": "修改库存成功!"})
	})

	router.POST("updateSku-compensate", func(c *gin.Context) {
		req := &proto.UpdateSkuReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		req.ProductSku.Stock += Number
		_, err := UpdateSkuClient.UpdateSku(context.TODO(), req)
		if err != nil {
			log.Println("/updateSku err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "回滚库存失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"updateSku-compensate": "SUCCESS", "Message": "回滚库存成功!"})
	})

	router.POST("/addCart", func(c *gin.Context) {
		req := &proto.AddCartReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		resp, err = AddCartClient.AddCart(context.TODO(), req)
		CartId = resp.ID
		//测试异常
		//err = errors.New("400", "测试异常", 400)
		if err != nil {
			log.Println("/addCart err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "新增购物车失败!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"addCart": "SUCCESS", "Message": "新增购物车成功!"})
	})
	router.POST("/addCart-compensate", func(c *gin.Context) {
		req := &proto.UpdateCartReq{}
		if err := c.BindJSON(req); err != nil {
			log.Fatalln(err)
		}
		req.Id = CartId
		resp, err := UpdateCartClient.UpdateCart(context.TODO(), req)
		CartId = resp.ID
		if err != nil {
			log.Println("/addCart-compensate err ", err)
			c.JSON(http.StatusOK, gin.H{"dtm_reslut": "FAILURE", "Message": "删除购物车失败!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"addCart-compensate": "SUCCESS", "Message": "删除购物车成功!"})
	})

	router.GET("/addShoppingCart", func(c *gin.Context) {
		number, _ := strconv.Atoi(c.Request.FormValue("number"))
		productId, _ := strconv.Atoi(c.Request.FormValue("productId"))
		productSkuId, _ := strconv.Atoi(c.Request.FormValue("productSkuId"))
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
		//拼接请求信息
		respErr := &proto.AddCartResp{}
		if err != nil || tokenResp.IsLogin == false {
			log.Println("GetUserToken  err : ", err)
			common.RespFail(c.Writer, "未登录！", respErr)
			return
		}
		log.Println("GetUserToken success : ", tokenResp)

		//拼接请求信息
		req := &proto.AddCartReq{
			Number:       int32(number),
			ProductId:    int32(productId),
			ProductSkuId: int32(productSkuId),
			UserId:       int32(sum),
			CreateUser:   int32(sum),
		}
		resp := &proto.AddCartResp{}
		//商品详情
		reqDetail := &proto.ProductDetailReq{
			Id: int32(productId),
		}
		respDetail, err := ProductDetailclient.ShowProductDetail(context.TODO(), reqDetail)
		if err != nil {
			log.Println("ShowProductDetail  err : ", err)
			common.RespFail(c.Writer, "查询商品详情失败！", respErr)
			return
		}
		if respDetail != nil {
			req.ProductName = respDetail.ProductDetail[0].Name
			req.ProductMainPicture = respDetail.ProductDetail[0].MainPicture
		}

		//log.Println(" /ShowProductDetail   resp   :", respDetail)
		//SKU详情
		reqDetail.Id = req.ProductSkuId
		respSkuDetail, err := ShowDetailSkuclient.ShowDetailSku(context.TODO(), reqDetail)
		//log.Println(" /ShowDetailSku   resp   :", respSkuDetail)
		//添加购物车  远程调用服务
		//log.Println(" /AddCart  req :", req)

		if respSkuDetail.ProductSku[0].Stock < req.Number {
			common.RespFail(c.Writer, "库存不足，添加失败", &proto.AddCartResp{})
			return
		}
		sku := respSkuDetail.ProductSku[0]
		sku.Stock -= req.Number
		Number = req.Number //
		updateSkuReq := &proto.UpdateSkuReq{
			ProductSku: sku,
		}
		resp.ProductSkuSimple = respSkuDetail.ProductSku[0]
		resp.ProductSimple = respDetail.ProductDetail[0]

		//全局事务
		gid := shortuuid.New()
		saga := dtmcli.NewSaga(DtmServer, gid).
			Add(QSBusi+"/updateSku", QSBusi+"/updateSku-compensate", updateSkuReq).
			Add(QSBusi+"/addCart", QSBusi+"/addCart-compensate", req)
		err = saga.Submit()
		if err != nil {
			log.Println("saga submit err :", err)
			common.RespFail(c.Writer, "添加失败", resp)
		}
		log.Println(" /saga submit   submit  :", gid)
		////writer  data  message  row  total  field
		common.RespOK(c.Writer, "请求成功", resp)
	})

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
	uuid := c.Request.Header["Uuid"][0]

	cc := common.GetInput(uuid)
	out := common.SQ(cc)
	sum := 0
	for o := range out {
		sum += o
	}

	req := &proto.AddCartReq{
		Number:       int32(number),
		ProductId:    int32(productId),
		ProductSkuId: int32(productSkuId),
		UserId:       int32(sum),
	}
	resp := &proto.AddCartResp{}

	// token 校验
	tokenReq := &proto.TokenReq{Uuid: uuid}
	tokenResp, err := GetUserTokenClient.GetUserToken(context.TODO(), tokenReq)
	if err != nil || tokenResp.IsLogin == false {
		log.Println("GetUserToken  err : ", err)
		common.RespFail(c.Writer, "未登录！", resp)
		return
	}
	log.Println("GetUserToken success : ", tokenResp)

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
	//SKU详情
	reqDetail.Id = req.ProductSkuId
	respSkuDetail, err := ShowDetailSkuclient.ShowDetailSku(context.TODO(), reqDetail)
	//添加购物车  远程调用服务
	if respSkuDetail.ProductSku[0].Stock < req.Number {
		common.RespFail(c.Writer, "库存不足，添加失败", &proto.AddCartResp{})
		return
	}

	sku := respSkuDetail.ProductSku[0]
	sku.Stock -= req.Number
	updateSkuReq := &proto.UpdateSkuReq{
		ProductSku: sku,
	}
	respUpdate, err := UpdateSkuClient.UpdateSku(context.TODO(), updateSkuReq)
	if err != nil {
		log.Println(" /UpdateSku  err :", err)
		common.RespFail(c.Writer, "修改库存失败！", resp)
		return
	}
	log.Println(" /UpdateSkuClient  resp :", respUpdate.IsSuccess)

	// 添加购物车，调用远程服务
	resp, err = AddCartClient.AddCart(context.TODO(), req)
	log.Println("/increase: ", resp)
	// 根据响应输出
	if err != nil {
		log.Println("addCart err ", err)
		updateSkuReq.ProductSku.Stock += req.Number
		_, err = UpdateSkuClient.UpdateSku(context.TODO(), updateSkuReq)
		log.Println("rollback sku  is Err :", err)
		common.RespFail(c.Writer, "添加购物车失败！", resp)
		return
	}
	resp.ProductSkuSimple = respSkuDetail.ProductSku[0]
	resp.ProductSimple = respDetail.ProductDetail[0]
	log.Println(" /AddCart  resp :", resp)
	common.RespOK(c.Writer, "增加购物车成功", resp)
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
