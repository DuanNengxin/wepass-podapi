package main

import (
	"fmt"
	"github.com/DuanNengxin/wepass-pod/common"
	"github.com/DuanNengxin/wepass-pod/plugin"
	pod "github.com/DuanNengxin/wepass-pod/proto"
	"github.com/DuanNengxin/wepass-podapi/handler"
	podapi "github.com/DuanNengxin/wepass-podapi/proto"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

var (
	//服务地址
	hostIp = "127.0.0.1"
	//服务地址
	serviceHost = hostIp
	//服务端口
	servicePort = "8082"
	//注册中心配置
	consulHost       = hostIp
	consulPort int64 = 8500
	//链路追踪
	tracerHost = hostIp
	tracerPort = 6831
	//熔断端口，每个服务不能重复
	hystrixPort = 9093
	//监控端口，每个服务不能重复
	//prometheusPort = 9192
)

func main() {

	consulReg := consul.NewRegistry(
		registry.Addrs(
			[]string{fmt.Sprintf("%s:%d", consulHost, consulPort)}...,
		),
	)

	tracer, i, err := common.NewTracer("wepass-podapi", fmt.Sprintf("%s:%d", tracerHost, tracerPort))
	if err != nil {
		return
	}
	defer i.Close()
	opentracing.SetGlobalTracer(tracer)

	// 添加熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	// 启动监听程序
	go func() {
		//http://192.168.0.112:9092/turbine/turbine.stream
		//看板访问地址 http://127.0.0.1:9002/hystrix，url后面一定要带 /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", strconv.Itoa(hystrixPort)), hystrixStreamHandler)
		fmt.Println("333")
		if err != nil {
			fmt.Println(err)
		}
	}()

	go common.PrometheusBoot()

	srv := micro.NewService(
		//自定义服务地址，必须要写在其它参数前面
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = serviceHost + ":" + servicePort
		})),
		micro.Name("wepass-podApi"),
		micro.Version("v1.0"),
		//指定服务端口
		micro.Address(":"+servicePort),
		//添加注册中心，
		micro.Registry(consulReg),
		//添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//作为客户端范围启动熔断
		micro.WrapClient(plugin.NewHystrixClientWrapper()),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
		//添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)
	// 初始化服务
	srv.Init()

	podService := pod.NewPodService("wepass-pod", srv.Client())
	// 创建服务句柄
	err = podapi.RegisterPodApiServiceHandler(srv.Server(), &handler.PodApi{PodService: podService})
	if err != nil {
		zap.S().Fatal(err)
	}
	if err := srv.Run(); err != nil {
		zap.S().Fatal(err)
	}

	select {}
}
