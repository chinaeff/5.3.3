package run

import (
	"context"
	geo2 "geotask_pprof/geo"
	cache2 "geotask_pprof/geo/cache"
	service2 "geotask_pprof/geo/module/courier/service"
	storage2 "geotask_pprof/geo/module/courier/storage"
	controller2 "geotask_pprof/geo/module/courierfacade/controller"
	service3 "geotask_pprof/geo/module/courierfacade/service"
	"geotask_pprof/geo/module/order/service"
	"geotask_pprof/geo/module/order/storage"
	router2 "geotask_pprof/geo/router"
	server2 "geotask_pprof/geo/server"
	order2 "geotask_pprof/geo/workers/order"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	// получение хоста и порта redis
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	// инициализация клиента redis

	rclient := cache2.NewRedisClient(host, port)

	// инициализация контекста с таймаутом
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// проверка доступности redis
	_, err := rclient.Ping().Result()
	if err != nil {
		return err
	}

	// инициализация разрешенной зоны
	allowedZone, _ := geo2.NewAllowedZone()
	// инициализация запрещенных зон
	disAllowedZone1, err := geo2.NewDisAllowedZone1()
	if err != nil {
		log.Fatal(err)
	}
	disAllowedZone2, err := geo2.NewDisAllowedZone2()
	if err != nil {
		log.Fatal(err)
	}
	disAllowedZones := []geo2.PolygonChecker{disAllowedZone1, disAllowedZone2}

	// инициализация хранилища заказов
	orderStorage := storage.NewOrderStorage(rclient)
	// инициализация сервиса заказов
	orderService := service.NewOrderService(orderStorage, allowedZone, disAllowedZones)

	orderGenerator := order2.NewOrderGenerator(orderService)
	orderGenerator.Run()

	oldOrderCleaner := order2.NewOrderCleaner(orderService)
	oldOrderCleaner.Run()

	// инициализация хранилища курьеров
	courierStorage := storage2.NewCourierStorage(rclient)
	// инициализация сервиса курьеров
	courierSevice := service2.NewCourierService(courierStorage, allowedZone, disAllowedZones)

	// инициализация фасада сервиса курьеров
	courierFacade := service3.NewCourierFacade(courierSevice, orderService)

	// инициализация контроллера курьеров
	courierController := controller2.NewCourierController(courierFacade)
	// инициализация роутера
	routes := router2.NewRouter(courierController)
	// инициализация сервера
	r := server2.NewHTTPServer()
	// инициализация группы роутов
	api := r.Group("/api")
	// инициализация роутов
	routes.CourierAPI(api)

	mainRoute := r.Group("/")

	routes.Swagger(mainRoute)
	// инициализация статических файлов
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))

	// запуск сервера
	//serverPort := os.Getenv("SERVER_PORT")

	if os.Getenv("ENV") == "prod" {
		certFile := "/app/certs/cert.pem"
		keyFile := "/app/certs/private.pem"
		return r.RunTLS(":443", certFile, keyFile)
	}

	return r.Run()
}
