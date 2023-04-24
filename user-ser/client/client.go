package main

import (
    "common"
    "context"
    "github.com/gin-gonic/gin"
    "github.com/go-micro/plugins/v4/registry/consul"
    "go-micro.dev/v4"
    "go-micro.dev/v4/registry"
    "go-micro.dev/v4/web"
    "log"
    "strconv"
    "user-ser/proto"
)

func getClient() proto.LoginService {
    // 注册到consul
    consulReg := consul.NewRegistry(func(options *registry.Options) {
        options.Addrs = []string{"101.34.10.3:8500"}
    })
    rpcServer := micro.NewService(
        micro.Registry(consulReg),
    )
    return proto.NewLoginService("shop-user", rpcServer.Client())
}

func UserLogin(c *gin.Context) {
    // 获取远程服务的客户端(获取服务)
    client := getClient()

    // 获取页面参数，拼接请求信息
    clientId, _ := strconv.Atoi(c.Request.FormValue("clientId"))
    phone := c.Request.FormValue("phone")
    systemId, _ := strconv.Atoi(c.Request.FormValue("systemId"))
    verificationCode := c.Request.FormValue("verificationCode")

    req := &proto.LoginRequest{
        ClientId:         int32(clientId),
        Phone:            phone,
        SystemId:         int32(systemId),
        VerificationCode: verificationCode,
    }

    // 远程调用服务
    resp, err := client.Login(context.TODO(), req)
    // 根据响应输出
    if err != nil {
        log.Println(err.Error())
        common.RespFail(c.Writer, "登录失败", resp)
        return
    }

    common.RespOK(c.Writer, "登录成功", resp)
}

func main() {
    router := gin.Default()

    router.Handle("GET", "toLogin", func(context *gin.Context) {
        context.String(200, "to login...")
    })

    router.GET("/login", UserLogin)

    service := web.NewService(
        web.Address(":6666"),
        web.Handler(router),
    )
    service.Run()
}
