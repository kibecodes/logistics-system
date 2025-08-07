package main

import (
	"log"
	"net/http"
	"os"

	"logistics-backend/handlers"
	"logistics-backend/internal/adapter/driveradapter"
	"logistics-backend/internal/adapter/inventoryadapter"
	"logistics-backend/internal/adapter/orderadapter"
	"logistics-backend/internal/repository/postgres"
	"logistics-backend/internal/router"
	deliveryUsecase "logistics-backend/internal/usecase/delivery"
	driverUsecase "logistics-backend/internal/usecase/driver"
	feedbackUsecase "logistics-backend/internal/usecase/feedback"
	inventoryUsecase "logistics-backend/internal/usecase/inventory"
	notificationUsecase "logistics-backend/internal/usecase/notification"
	orderUsecase "logistics-backend/internal/usecase/order"
	paymentUsecase "logistics-backend/internal/usecase/payment"
	userUsecase "logistics-backend/internal/usecase/user"

	_ "logistics-backend/docs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title Logistics API
// @version 1.0
// @description This is the API for logistics operations.
// @host localhost:8000
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL not set")
	}

	publicApiBaseUrl := os.Getenv("PUBLIC_API_BASE_URL")
	if publicApiBaseUrl == "" {
		log.Fatal("PUBLIC_API_BASE_URL not set")
	}

	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Set up repositories
	userRepo := postgres.NewUserRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	driverRepo := postgres.NewDriverRepository(db)
	deliveryRepo := postgres.NewDeliveryRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	feedbackRepo := postgres.NewFeedbackRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	inventoryRepo := postgres.NewInventoryRespository(db)

	// Set up usecase
	dUsecase := driverUsecase.NewUseCase(driverRepo)
	driveradapter := driveradapter.DriverUseCaseAdapter{UseCase: dUsecase}
	uUsecase := userUsecase.NewUseCase(userRepo, &driveradapter)
	pUsecase := paymentUsecase.NewUseCase(paymentRepo)
	fUsecase := feedbackUsecase.NewUseCase(feedbackRepo)
	nUsecase := notificationUsecase.NewUseCase(notificationRepo)
	iUsecase := inventoryUsecase.NewUseCase(inventoryRepo)
	inventoryAdapter := inventoryadapter.InventoryUseCaseAdapter{UseCase: iUsecase}
	oUsecase := orderUsecase.NewUseCase(orderRepo, &inventoryAdapter)
	orderAdapter := orderadapter.OrderUseCaseAdapter{UseCase: oUsecase}
	eUsecase := deliveryUsecase.NewUseCase(deliveryRepo, &orderAdapter)

	// Set up Handlers
	userHandler := handlers.NewUserHandler(uUsecase)
	orderHandler := handlers.NewOrderHandler(oUsecase)
	driverHandler := handlers.NewDriverHandler(dUsecase)
	deliveryHandler := handlers.NewDeliveryHandler(eUsecase)
	paymentHandler := handlers.NewPaymentHandler(pUsecase)
	feedbackHandler := handlers.NewFeedbackHandler(fUsecase)
	notificationHandler := handlers.NewNotificationHandler(nUsecase)
	inventoryHandler := handlers.NewInventoryHandler(iUsecase)

	// Start server
	r := router.NewRouter(userHandler, orderHandler, driverHandler, deliveryHandler, paymentHandler, feedbackHandler, notificationHandler, inventoryHandler, publicApiBaseUrl)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("could not start server at: %v", err)
	}
}
