package main

import (
    "common"
    "github.com/go-micro/plugins/v4/registry/consul"
    opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
    "github.com/opentracing/opentracing-go"
    "go-micro.dev/v4"
    "go-micro.dev/v4/registry"
    "log"
    "product-ser/domain/repository"
    "product-ser/domain/service"
    "product-ser/handler"
    "product-ser/proto"
    "time"
)

const (
    consulStr      = "http://101.34.10.3:8500"
    consulReistStr = "101.34.10.3:8500"
    fileKey        = "mysql-product"
)

func main() {
    // 配置中心
    consulConfig, err := common.GetConsulConfig(consulStr, fileKey)
    if err != nil {
        log.Fatalf("consul config error: %+v", err)
    }

    // 注册中心
    consulReist := consul.NewRegistry(func(options *registry.Options) {
        options.Addrs = []string{consulReistStr}
    })

    //链路追踪实例化
    trancer, io, err := common.NewTraner("shop-product", "101.34.10.3:6831")
    if err != nil {
        log.Fatalln(err)
    }
    defer io.Close()
    // 设置全局的 tracing
    opentracing.SetGlobalTracer(trancer)

    rpcService := micro.NewService(
        micro.RegisterTTL(time.Second*30),
        micro.RegisterInterval(time.Second*30),
        micro.Name("shop-product"),
        micro.Address(":8082"),
        micro.Version("v1"),
        micro.Registry(consulReist),
        micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
    )

    // 初始化DB
    db, err := common.GetMysqlFromConsul(consulConfig)
    if err != nil {
        log.Fatalf("mysql error: %+v", err)
    }

    // 创建服务实例
    productDataService := service.NewProductDataService(repository.NewProductRepository(db))

    // 注册handler
    err = proto.RegisterPageHandler(rpcService.Server(), &handler.ProductHandler{ProductDataService: productDataService})
    if err != nil {
        log.Fatalln("register handler error: ", err)
    }

    err = proto.RegisterShowProductDetailHandler(rpcService.Server(), &handler.ProductHandler{ProductDataService: productDataService})
    if err != nil {
        log.Fatalln("register handler error: ", err)
    }

    err = proto.RegisterShowProductSkuHandler(rpcService.Server(), &handler.ProductHandler{ProductDataService: productDataService})
    if err != nil {
        log.Fatalln("register handler error: ", err)
    }

    // 启动服务
    if err := rpcService.Run(); err != nil {
        log.Fatalln("start user service error: ", err)
    }
}
