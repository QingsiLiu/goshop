package main

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
	"time"
	"user-ser/common"
	"user-ser/domain/repository"
	"user-ser/domain/service"
	"user-ser/handler"
	"user-ser/proto"
)

const (
	consulStr = "http://101.34.10.3:8500"
	fileKey   = "mysql-user"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig(consulStr, fileKey)
	if err != nil {
		log.Panicf("consul config error: %+v", err)
	}

	// consul 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{consulStr}
	})
	rpcService := micro.NewService(
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*30),
		micro.Name("shop-user"),
		micro.Address(":8081"),
		micro.Version("v1"),
		micro.Registry(consulRegistry),
	)

	// 初始化DB
	db, err := common.GetMysqlFromConsul(consulConfig)
	if err != nil {
		log.Panicf("mysql error: %+v", err)
	}

	// 创建服务实例
	userDataService := service.NewUserDataService(repository.NewUserRepository(db))

	// 注册handler
	err = proto.RegisterLoginHandler(rpcService.Server(), &handler.User{UserDataService: userDataService})
	if err != nil {
		log.Panicln("register handler error: ", err)
	}

	// 启动服务
	if err := rpcService.Run(); err != nil {
		log.Panicln("start user service error: ", err)
	}
}
