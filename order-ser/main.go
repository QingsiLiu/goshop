package main

import (
	"common"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	ratelimiter "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
	"order-ser/domain/repository"
	"order-ser/domain/service"
	"order-ser/handler"
	"order-ser/proto"
	"time"
)

const (
	consulStr      = "http://101.34.10.3:8500"
	consulReistStr = "101.34.10.3:8500"
	fileKey        = "mysql-product"
	QPS            = 100
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
	trancer, io, err := common.NewTraner("shop-order", "101.34.10.3:6831")
	if err != nil {
		log.Fatalln(err)
	}
	defer io.Close()
	// 设置全局的 tracing
	opentracing.SetGlobalTracer(trancer)

	// 开始监控 Prometheus
	common.PrometheusBoot(9092)

	rpcService := micro.NewService(
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*30),
		micro.Name("shop-order"),
		micro.Address(":8084"),
		micro.Version("v1"),
		micro.Registry(consulReist),
		// 链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(QPS)),
		// 添加监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 初始化DB
	db, err := common.GetMysqlFromConsul(consulConfig)
	if err != nil {
		log.Fatalf("mysql error: %+v", err)
	}

	// 创建服务实例
	orderService := service.NewOrderService(repository.NewOrderRepository(db))

	// 注册handler
	err = proto.RegisterFindOrderHandler(rpcService.Server(), &handler.OrderHandler{OrderService: orderService})
	if err != nil {
		log.Fatalln("register handler error: ", err)
	}
	err = proto.RegisterAddTradeOrderHandler(rpcService.Server(), &handler.OrderHandler{OrderService: orderService})
	if err != nil {
		log.Fatalln("register handler error: ", err)
	}
	err = proto.RegisterUpdateTradeOrderHandler(rpcService.Server(), &handler.OrderHandler{OrderService: orderService})
	if err != nil {
		log.Fatalln("register handler error: ", err)
	}

	// 启动服务
	if err := rpcService.Run(); err != nil {
		log.Fatalln("start user service error: ", err)
	}
}
