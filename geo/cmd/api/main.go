package main

import (
	"fmt"
	"geotask_pprof/geo/module/courier/models"
	"geotask_pprof/proxy"
	"geotask_pprof/rbmq"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"net/rpc"
	"os"
	"time"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Request total",
		},
		[]string{"endpoint"},
	)

	cacheDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_duration_seconds",
			Help:    "Cache duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	dbDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_duration_seconds",
			Help:    "DB duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	externalAPIDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_duration_seconds",
			Help:    "API duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// @title courier service
// @version 1.0
// @description courier service
// @host localhost:8080
// @BasePath /api/v1
//
//go:generate swagger generate spec -o ../public/swagger.json --scan-models
func main() {
	go rbmq.StartNotificationService()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	messageBroker := os.Getenv("MESSAGE_BROKER")

	var _ queue.MessageQueuer
	switch messageBroker {
	case "RabbitMQ":
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Fatal(err)
		}
		_, err = rabbitmq.NewRabbitMQ(conn)
		if err != nil {
			log.Fatal(err)
		}
	case "Kafka":
		conn, err := kafka.Dial("tcp", "localhost:9092")
		if err != nil {
			log.Fatal(err)
		}
		_, err = kafka.NewKafka(conn)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unknown message broker specified.")
	}

	protocol := os.Getenv("RPC_PROTOCOL")

	provider := NewGeoProvider()

	geoService := &GeoService{Provider: provider}

	rpc.Register(geoService)

	switch protocol {
	case "json-rpc":
		startJSONRPCServer()
	case "rpc":
		startRPCServer()
	default:
		fmt.Println("Unknown RPC protocol specified.")
	}

	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(cacheDuration)
	prometheus.MustRegister(dbDuration)
	prometheus.MustRegister(externalAPIDuration)

	http.Handle("/metrics", promhttp.Handler())

	log.Println(http.ListenAndServe("localhost:9090", nil))
	go func() {
		http.Handle("/mycustompath/pprof/", http.HandlerFunc(pprofHandler))
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	proxyConfig := proxy.ReverseProxyConfig{
		BackendURLs: []string{"http://geo1:8080", "http://geo2:8080", "http://geo3:8080"},
	}

	proxyN := proxy.NewReverseProxy(proxyConfig)

	router := gin.Default()
	router.Use()

	rpcServer := rpc.NewServer()

	_, err = net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	if err := rpcServer.Accept; err != nil {
		log.Fatal(err)
	}

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
	))

	router.Any("/api/address/*any", func(c *gin.Context) {
		proxyN.ServeHTTP(c.Writer, c.Request)
	})
	router.POST("/move-courier", func(c *gin.Context) {
		startTime := time.Now()
		var courierLocation models.Point
		if err := c.BindJSON(&courierLocation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Courier location updated"})
		duration := time.Since(startTime).Seconds()

		requestDuration.WithLabelValues("/move-courier").Observe(duration)
		requestCount.WithLabelValues("/move-courier").Inc()
	})
	router.POST("/api/address/*any", rbmq.RbHandler)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("error conn to rabbit", err)
	}
	defer conn.Close()

	rbmq.RabbitMQ, err = rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		log.Fatal("error creating rabbit", err)
	}

	router.Run(":8080")
	if err := saveProfile("/mycustompath/pprof/profile", "profile.pprof"); err != nil {
		fmt.Println("Error saving profile:", err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != "some token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func pprofHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profileType := r.URL.Path[len("/mycustompath/pprof/"):]
	switch profileType {
	case "allocs", "block", "cmdline", "goroutine", "heap", "mutex", "profile", "threadcreate", "trace":
		pprof.Index(w, r)
	default:
		http.NotFound(w, r)
	}
}
func isAuthorized(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	return token == "some token"
}
func saveProfile(endpoint, filename string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:6060%s", endpoint), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

//func rbHandler(c *gin.Context) {
//	if !rateLimiter.Take().IsZero() {
//		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
//		return
//	}
//
//	err := rabbitMQ.Publish("user_rate_limit_exceeded", []byte("User rate limit exceeded"))
//	if err != nil {
//		log.Fatal("error public reabbit", err)
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Request processed successfully"})
//}
//
//func startNotificationService() {
//	messages, err := rabbitMQ.Subscribe("user_rate_limit_exceeded")
//	if err != nil {
//		log.Fatal("Ошибка подписки на канал сообщений RabbitMQ:", err)
//	}
//
//	processMessages(messages)
//}
//
//func processMessages(messages <-chan queue.Message) {
//	for msg := range messages {
//		msgString := msg
//
//		sendEmailNotification(msgString)
//		sendSMSNotification(msgString)
//	}
//}
//
//func sendEmailNotification(msg queue.Message) {
//	log.Println("Email:", msg)
//}
//
//func sendSMSNotification(msg queue.Message) {
//	log.Println("SMS:", msg)
//}
