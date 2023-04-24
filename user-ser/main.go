package main

import (
    "common"
    "github.com/go-micro/plugins/v4/registry/consul"
    "go-micro.dev/v4"
    "go-micro.dev/v4/registry"
    "log"
    "time"
    "user-ser/domain/repository"
    "user-ser/domain/service"
    "user-ser/handler"
    "user-ser/proto"
)

const (
    consulStr      = "http://101.34.10.3:8500"
    consulReistStr = "101.34.10.3:8500"
    fileKey        = "mysql-user"
    redisKwy       = "redis"
)

func main() {
    // 配置中心
    consulConfig, err := common.GetConsulConfig(consulStr, fileKey)
    if err != nil {
        log.Panicf("consul config error: %+v", err)
    }

    // consul 注册中心
    consulRegistry := consul.NewRegistry(func(options *registry.Options) {
        options.Addrs = []string{consulReistStr}
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
    // 初始化redis
    consulRedisConfig, err := common.GetConsulConfig(consulStr, redisKwy)
    if err != nil {
        log.Panicf("consul config error: %+v", err)
    }
    red, _ := common.GetRedisFromConsul(consulRedisConfig)
    common.SetUserToken(red, "uuid111", []byte("tokenxxx"), time.Duration(1)*time.Hour)
    res, _ := common.GetUserToken(red, "uuid111")
    log.Println(res)

    // 创建服务实例
    userDataService := service.NewUserDataService(repository.NewUserRepository(db, red))

    // 注册handler
    err = proto.RegisterLoginHandler(rpcService.Server(), &handler.UserHandler{UserDataService: userDataService})
    if err != nil {
        log.Panicln("register handler error: ", err)
    }

    err = proto.RegisterGetUserTokenHandler(rpcService.Server(), &handler.UserHandler{UserDataService: userDataService})
    if err != nil {
        log.Panicln("register handler error: ", err)
    }

    // 启动服务
    if err := rpcService.Run(); err != nil {
        log.Panicln("start user service error: ", err)
    }
}
