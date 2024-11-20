package main

import (
	"github.com/llchhh/spektr-account-api/auth"
	"github.com/llchhh/spektr-account-api/internal/repository/api"
	"github.com/llchhh/spektr-account-api/internal/rest"
	"github.com/llchhh/spektr-account-api/internal/rest/middleware"
	"github.com/llchhh/spektr-account-api/profile"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
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

	//// Prepare Repository
	//authorRepo := mysqlRepo.NewAuthorRepository(dbConn)
	//articleRepo := mysqlRepo.NewArticleRepository(dbConn)
	//
	//// Build service Layer
	//svc := article.NewService(articleRepo, authorRepo)
	//rest.NewArticleHandler(e, svc)

	authRepo := api.NewAuthRepository(os.Getenv("BASE_URL"))
	authSvc := auth.NewService(authRepo)
	rest.NewAuthHandler(e, authSvc)
	profileRepo := api.NewProfileRepository(os.Getenv("BASE_URL"))
	profileSvc := profile.NewService(profileRepo)
	rest.NewProfileHandler(e, profileSvc)
	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	log.Fatal(e.Start(address)) //nolint
}
