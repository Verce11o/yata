package app

import (
	"fmt"
	"github.com/Verce11o/yata/internal/config"
	"github.com/Verce11o/yata/internal/http"
	"github.com/Verce11o/yata/internal/http/auth"
	"github.com/Verce11o/yata/internal/lib/logger"
	"github.com/Verce11o/yata/internal/lib/response"
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

	// Init services
	services := service.NewServices(cfg)

	// Init handlers
	authHandler := auth.NewHandler(log, services, validator)
	handlers := http.NewHandlers(authHandler)

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
