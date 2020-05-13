package main

import (
	"github.com/micro/go-config"
	//"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"net/http"

	//"github.com/opentracing/opentracing-go"
	//"os"
	"shopping/order/handler"
	"shopping/order/model"
	order "shopping/order/proto/order"
	"shopping/order/repository"
	product "shopping/product/proto/product"

	"github.com/micro/go-plugins/broker/rabbitmq"

	//"go.opencensus.io/trace"
	//"go.opencensus.io/exporter/zipkin"
	//wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opencensus"
	//openzipkin "github.com/openzipkin/zipkin-go"
	//zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"

	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing"
	//"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	wrapperPrometheus "github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	//加载配置项
	err := config.LoadFile("./config.json")
	if err != nil {
		log.Fatalf("Could not load config file: %s", err.Error())
		return
	}
	conf := config.Map()

	//db
	db, err := CreateConnection(conf["mysql"].(map[string]interface{}))
	defer db.Close()

	db.AutoMigrate(&model.Order{})

	if err != nil {
		log.Fatalf("connection error : %v \n", err)
	}

	repo := &repository.Order{db}

	//broker
	b := rabbitmq.NewBroker(
		broker.Addrs(config.Get("rabbitmq_addr").String("")),
	)

	// boot trace
	TraceBoot()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.order"),
		micro.Version("latest"),
		micro.Broker(b),
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper()),
		micro.WrapClient(wrapperTrace.NewClientWrapper()),
		micro.WrapHandler(wrapperPrometheus.NewHandlerWrapper()),
	)

	// boot prometheus
	PrometheusBoot()

	// Initialise service
	service.Init()

	//创建消息发布者
	publisher := micro.NewPublisher("notification.submit", service.Client())

	//product-srv client
	productCli := product.NewProductService("go.micro.srv.product", service.Client())

	// Register Handler
	order.RegisterOrderServiceHandler(service.Server(), &handler.Order{repo, productCli, publisher})

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.order", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.order", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

//trace opencensus+zipkin
//func TraceBoot() {
//	apiURL := "http://192.168.0.111:9411/api/v2/spans"
//	hostPort,_ := os.Hostname()
//	serviceName := "go.micro.srv.order"
//
//	localEndpoint, err := openzipkin.NewEndpoint(serviceName, hostPort)
//	if err != nil {
//		log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
//	}
//	reporter := zipkinHTTP.NewReporter(apiURL)
//	ze := zipkin.NewExporter(reporter, localEndpoint)
//	trace.RegisterExporter(ze)
//	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
//	return
//}

//trace opentracing+zipkin
//func TraceBoot() {
//	apiURL := "http://192.168.0.111:9411/api/v1/spans"
//	hostPort,_ := os.Hostname()
//	serviceName := "go.micro.srv.order"
//
//	collector, err := zipkin.NewHTTPCollector(apiURL)
//	if err != nil {
//		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
//		return
//	}
//	tracer, err := zipkin.NewTracer(
//		zipkin.NewRecorder(collector, false, hostPort, serviceName),
//	)
//	if err != nil {
//		log.Fatalf("unable to create Zipkin tracer: %v", err)
//		return
//	}
//	opentracing.InitGlobalTracer(tracer)
//	return
//}


//trace opentracing+Jaeger
func TraceBoot() {
	serviceName := "go.micro.srv.order"
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:  "192.168.0.111:9412",
		},
	}

	closer, err := cfg.InitGlobalTracer(
		serviceName,
	)
	if err != nil {
		log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()

	//opentracing.InitGlobalTracer(tracer)
	return
}


func PrometheusBoot(){
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("192.168.0.110:8085", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}