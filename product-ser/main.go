package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
	"product-ser/common"
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
		log.Panicf("consul config error: %+v", err)
	}

	// 注册中心
	consulReist := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{consulReistStr}
	})
	rpcService := micro.NewService(
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*30),
		micro.Name("shop-product"),
		micro.Address(":8082"),
		micro.Version("v1"),
		micro.Registry(consulReist),
	)

	// 初始化DB
	db, err := common.GetMysqlFromConsul(consulConfig)
	if err != nil {
		log.Panicf("mysql error: %+v", err)
	}

	// 创建服务实例
	productDataService := service.NewProductDataService(repository.NewProductRepository(db))

	// 注册handler
	err = proto.RegisterPageHandler(rpcService.Server(), &handler.ProductHandler{ProductDataService: productDataService})
	if err != nil {
		log.Panicln("register handler error: ", err)
	}

	// 启动服务
	if err := rpcService.Run(); err != nil {
		log.Panicln("start user service error: ", err)
	}
}
