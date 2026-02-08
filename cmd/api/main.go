package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidcm146/assets-management-be.git/internal/config"
	"github.com/davidcm146/assets-management-be.git/internal/database"
	"github.com/davidcm146/assets-management-be.git/internal/handler"
	"github.com/davidcm146/assets-management-be.git/internal/infrastructure/cloudinary"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"github.com/davidcm146/assets-management-be.git/internal/router"
	"github.com/davidcm146/assets-management-be.git/internal/server"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/davidcm146/assets-management-be.git/internal/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := validator.RegisterValidators(); err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.NewDB(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cloudCfg := config.LoadCloudinaryConfig()
	cld, _ := cloudinary.NewCloudinary(cloudCfg)
	uploader := cloudinary.NewCloudinaryUploader(cld)

	userRepo := repository.NewUserRepository(db)
	loanSlipRepo := repository.NewLoanSlipRepository(db)

	authService := service.NewAuthService(userRepo)
	loanSlipService := service.NewLoanSlipService(loanSlipRepo, uploader)

	authHandler := handler.NewAuthHandler(authService, userRepo)
	loanSlipHandler := handler.NewLoanSlipHandler(loanSlipService, uploader)

	engine := gin.New()
	r := router.NewRouter(router.RouterParams{
		Engine: engine,
		Handlers: &router.Handlers{
			AuthHandler:     authHandler,
			LoanSlipHandler: loanSlipHandler,
		},
	})
	srv := server.NewServer(r, cfg.Server.Port)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(context.Background())
}
