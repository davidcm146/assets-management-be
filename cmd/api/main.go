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
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"github.com/davidcm146/assets-management-be.git/internal/router"
	"github.com/davidcm146/assets-management-be.git/internal/server"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
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

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	engine := gin.New()
	r := router.NewRouter(router.RouterParams{
		Engine: engine,
		Handlers: &router.Handlers{
			AuthHandler: authHandler,
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
