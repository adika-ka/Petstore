package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"petstore/infrastructure"
	"petstore/internal/config"
	"petstore/internal/controller"
	"petstore/internal/db"
	"petstore/internal/middleware"
	"petstore/internal/repository"
	"petstore/internal/service"
	"syscall"
	"time"

	_ "petstore/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

// @title           Petstore API
// @version         1.0
// @description     This is a sample Petstore server.
// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.InitJWT()
	dbConn, err := db.InitDBAndMigrate()
	if err != nil {
		log.Fatalf("Failed to init DB and run migrations: %v", err)
	}
	defer dbConn.Close()

	petRepo := repository.NewPetRepository(dbConn)
	petService := service.NewPetService(petRepo)
	responder := infrastructure.NewJSONResponder()

	petController := &controller.PetController{
		Service:   petService,
		Responder: responder,
	}

	orderRepo := repository.NewOrderRepository(dbConn)
	orderService := service.NewOrderService(orderRepo)

	orderController := &controller.OrderController{
		Service:   orderService,
		Responder: responder,
	}

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)

	userController := &controller.UserController{
		Service:   userService,
		Responder: responder,
	}

	r := chi.NewRouter()

	controller.RegisterUserRoutes(r, userController)
	controller.RegisterOrderRoutes(r, orderController)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuthMiddleware)
		controller.RegisterPetRoutes(protected, petController)
		protected.Get("/store/inventory", controller.GetInventory(orderController))
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server is starting at :8080")
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
