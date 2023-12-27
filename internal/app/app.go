package app

import (
	"fmt"
	"github.com/Verce11o/yata/internal/config"
	"github.com/Verce11o/yata/internal/http"
	"github.com/Verce11o/yata/internal/http/auth"
	"github.com/Verce11o/yata/internal/http/comments"
	"github.com/Verce11o/yata/internal/http/middleware"
	"github.com/Verce11o/yata/internal/http/notifications"
	"github.com/Verce11o/yata/internal/http/tweets"
	"github.com/Verce11o/yata/internal/http/websocket"
	"github.com/Verce11o/yata/internal/lib/logger"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
	"github.com/Verce11o/yata/internal/lib/response"
	"github.com/Verce11o/yata/internal/rabbitmq"
	"github.com/Verce11o/yata/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	app := fiber.New()
	app.Use(cors.New())

	log := logger.NewLogger(cfg.Mode)
	validator := response.NewValidator()

	// Init metrics
	tracer := trace.InitTracer("http")

	// Init services
	services := service.NewServices(cfg, tracer)

	// Init websocket
	websocketHandler := websocket.NewHandler(log, tracer.Tracer, services)

	// Init broker
	amqpConn := rabbitmq.NewAmqpConnection(cfg.RabbitMQ)
	notificationConsumer := rabbitmq.NewNotificationConsumer(amqpConn, log, tracer.Tracer)

	go func() {
		err := notificationConsumer.StartConsumer(
			cfg.RabbitMQ.QueueName,
			cfg.RabbitMQ.ConsumerTag,
			cfg.RabbitMQ.ExchangeName,
			cfg.RabbitMQ.BindingKey,
			websocketHandler.Clients,
		)

		if err != nil {
			log.Errorf("StartConsumerErr: %v", err.Error())
		}

	}()

	// Init middleware
	middlewareHandler := middleware.NewMiddlewareHandler(log, tracer.Tracer, services, cfg, validator)

	// Init handlers
	authHandler := auth.NewHandler(log, tracer.Tracer, services, validator)
	tweetHandler := tweets.NewHandler(log, tracer.Tracer, services, validator)
	commentHandler := comments.NewHandler(log, tracer.Tracer, services, validator)
	notificationHandler := notifications.NewHandler(log, tracer.Tracer, services, validator)

	handlers := http.NewHandlers(authHandler, tweetHandler, commentHandler, notificationHandler, websocketHandler, middlewareHandler)

	handlers.InitRoutes(app)

	app.Use(fiberLogger.New())

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.HTTPServer.Port)); err != nil {
			log.Fatal("error while running server: ", err)
		}
	}()

	log.Infof("Server is running on port: %v", cfg.HTTPServer.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Server exiting..")

	if err := app.Shutdown(); err != nil {
		log.Fatal("Server Shutdown error: ", err)
	}

}
