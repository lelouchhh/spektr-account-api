// @title Example API
// @version 1.0
// @description This is a sample server for managing users.
// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/llchhh/spektr-account-api/auth"
	_ "github.com/llchhh/spektr-account-api/docs" // Import generated docs
	"github.com/llchhh/spektr-account-api/internal/repository/api"
	"github.com/llchhh/spektr-account-api/internal/rest"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
	"github.com/llchhh/spektr-account-api/notification"
	"github.com/llchhh/spektr-account-api/profile"
	"github.com/llchhh/spektr-account-api/repair"
	"github.com/swaggo/http-swagger" // Swagger UI handler
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// prepare echo
	e := echo.New()
	e.Use(middleware.CORS)

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Prepare Repositories
	authRepo := api.NewAuthRepository(os.Getenv("BASE_URL"))
	authSvc := auth.NewService(authRepo)
	rest.NewAuthHandler(e, authSvc)

	profileRepo := api.NewProfileRepository(os.Getenv("BASE_URL"))
	profileSvc := profile.NewService(profileRepo)
	rest.NewProfileHandler(e, profileSvc)

	notiRepo := api.NewNotificationRepository(os.Getenv("BASE_URL"))
	notiSvc := notification.NewService(notiRepo)
	rest.NewNotificationHandler(e, notiSvc)

	repairRepo := api.NewRepairRepository(os.Getenv("BASE_URL"))
	repairSvc := repair.NewService(repairRepo)
	rest.NewRepairHandler(e, *repairSvc)

	// Получаем API ключ из переменной окружения
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment variables")
	}

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	// Указываем путь до сертификатов
	certFile := "/etc/letsencrypt/live/www.969975-cv27771.tmweb.ru/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/www.969975-cv27771.tmweb.ru/privkey.pem"

	// Проверка наличия файлов сертификата и ключа
	_, certErr := os.Stat(certFile)
	_, keyErr := os.Stat(keyFile)

	// Если сертификат или ключ отсутствуют, запускаем сервер без SSL
	if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
		log.Println("SSL certificates not found, starting HTTP server instead of HTTPS.")
		e.GET("/swagger/*", echo.WrapHandler(httpSwagger.WrapHandler))

		swaggerGroup := e.Group("/swagger")
		swaggerGroup.Use(middleware.APIKey(apiKey)) // Используем middleware для проверки API ключа

		swaggerGroup.GET("/*", echo.WrapHandler(httpSwagger.WrapHandler))
		log.Fatal(e.Start(address)) // Start HTTP server without SSL
	} else {
		log.Printf("SSL certificates found, starting HTTPS server on %s...\n", address)
		e.GET("/swagger/*", echo.WrapHandler(httpSwagger.WrapHandler))

		swaggerGroup := e.Group("/swagger")
		swaggerGroup.Use(middleware.APIKey(apiKey)) // Используем middleware для проверки API ключа

		swaggerGroup.GET("/*", echo.WrapHandler(httpSwagger.WrapHandler))
		log.Fatal(e.StartTLS(address, certFile, keyFile)) // Start HTTPS server
	}
}
