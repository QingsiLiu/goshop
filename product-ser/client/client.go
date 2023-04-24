package main

import (
    "common"
    "context"
    "github.com/gin-gonic/gin"
    "github.com/go-micro/plugins/v4/registry/consul"
    opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
    "github.com/opentracing/opentracing-go"
    "go-micro.dev/v4"
    "go-micro.dev/v4/registry"
    "go-micro.dev/v4/web"
    "log"
    "product-ser/proto"
    "strconv"
)

var (
    client    proto.PageService
    clientA   proto.ShowProductDetailService
    clientSKU proto.ShowProductSkuService
)

func main() {
    router := gin.Default()

    // 注册到consul
    consulReg := consul.NewRegistry(func(options *registry.Options) {
        options.Addrs = []string{"101.34.10.3:8500"}
    })

    // 初始化链路追踪的jaeger
    trancer, io, err := common.NewTraner("shop-product-client", "101.34.10.3:6831")
    if err != nil {
        log.Fatalln(err)
    }
    defer io.Close()
    opentracing.SetGlobalTracer(trancer)

    rpcServer := micro.NewService(
        micro.Registry(consulReg), // 服务发现
        micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
    )

    client = proto.NewPageService("shop-product", rpcServer.Client())
    clientA = proto.NewShowProductDetailService("shop-product", rpcServer.Client())
    clientSKU = proto.NewShowProductSkuService("shop-product", rpcServer.Client())

    // 具体的路由
    // 分页查询
    router.GET("/page", ProductPage)
    // 查询商品详情
    router.GET("/showProductDetail", ShowProductDetail)
    // 查询商品SKU
    router.GET("/sku", ShowProductSKU)

    service := web.NewService(
        web.Address(":6667"),
        web.Name("shop-product-client"),
        web.Registry(consulReg),
        web.Handler(router),
    )
    service.Run()
}

func ProductPage(c *gin.Context) {
    // 获取页面参数，拼接请求信息
    length, _ := strconv.Atoi(c.Request.FormValue("length"))
    pageIndex, _ := strconv.Atoi(c.Request.FormValue("pageIndex"))

    req := &proto.PageReq{
        Length:    int32(length),
        PageIndex: int32(pageIndex),
    }

    // 远程调用服务
    resp, err := client.Page(context.TODO(), req)
    // 根据响应输出
    if err != nil {
        log.Println(err.Error())
        common.RespFail(c.Writer, "获得商品列表失败", resp)
        return
    }

    common.RespListOK(c.Writer, "获得商品列表成功", resp, resp.Rows, resp.Total, "")
}

func ShowProductDetail(c *gin.Context) {
    //获取页面参数
    id, _ := strconv.Atoi(c.Request.FormValue("id"))
    //拼接请求信息
    req := &proto.ProductDetailReq{
        Id: int32(id),
    }
    //远程调用服务
    resp, err := clientA.ShowProductDetail(context.TODO(), req)
    log.Println(" /showProductDetail  :", resp)
    //根据响应做输出
    if err != nil {
        log.Println(err.Error())
        //c.String(http.StatusBadRequest, "search failed !")
        common.RespFail(c.Writer, "请求失败", resp)
        return
    }
    ////writer  data  message  row  total  field
    common.RespOK(c.Writer, "请求成功", resp)
}

func ShowProductSKU(c *gin.Context) {
    //获取远程服务的客户端 client
    //获取页面参数
    id, _ := strconv.Atoi(c.Request.FormValue("productId"))
    //拼接请求信息
    req := &proto.ProductSkuReq{
        ProductId: int32(id),
    }
    //远程调用服务
    resp, err := clientSKU.ShowProductSku(context.TODO(), req)
    log.Println(" /sku  :", resp)
    //根据响应做输出
    if err != nil {
        log.Println(err.Error())
        //c.String(http.StatusBadRequest, "search failed !")
        common.RespFail(c.Writer, "请求失败", resp)
        return
    }
    ////writer  data  message  row  total  field
    common.RespListOK(c.Writer, "请求成功", resp, 0, 0, "请求成功")

}
