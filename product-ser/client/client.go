package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"log"
	"product-ser/common"
	"product-ser/proto"
	"strconv"
)

func getClient() proto.PageService {
	// 注册到consul
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"101.34.10.3:8500"}
	})
	rpcServer := micro.NewService(
		micro.Registry(consulReg),
	)
	return proto.NewPageService("shop-product", rpcServer.Client())
}

func ProductPage(c *gin.Context) {
	// 获取远程服务的客户端(获取服务)
	client := getClient()

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

	common.RespOK(c.Writer, "获得商品列表成功", resp)
}

func main() {
	router := gin.Default()

	router.GET("/page", ProductPage)

	service := web.NewService(
		web.Address(":6667"),
		web.Handler(router),
	)
	service.Run()
}
