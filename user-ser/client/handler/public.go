package handler

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
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
