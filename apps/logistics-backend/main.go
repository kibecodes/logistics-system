package main

import (
	"log"
	"net/http"
	"os"

	"logistics-backend/handlers"
	deliveryadapter "logistics-backend/internal/adapters/delivery"
	driveradapter "logistics-backend/internal/adapters/driver"
	inventoryadapter "logistics-backend/internal/adapters/inventory"
	orderadapter "logistics-backend/internal/adapters/order"
	useradapter "logistics-backend/internal/adapters/user"
	"logistics-backend/internal/repository/postgres"
	"logistics-backend/internal/router"
	deliveryUsecase "logistics-backend/internal/usecase/delivery"
	driverUsecase "logistics-backend/internal/usecase/driver"
	feedbackUsecase "logistics-backend/internal/usecase/feedback"
	inventoryUsecase "logistics-backend/internal/usecase/inventory"
	inviteUsecase "logistics-backend/internal/usecase/invite"
	notificationUsecase "logistics-backend/internal/usecase/notification"
	orderUsecase "logistics-backend/internal/usecase/order"
	paymentUsecase "logistics-backend/internal/usecase/payment"
	userUsecase "logistics-backend/internal/usecase/user"

	"logistics-backend/internal/application"

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
	inviteRepo := postgres.NewInviteRepository(db)

	// Set up usecase
	// Individual
	inviteUC := inviteUsecase.NewUseCase(inviteRepo)
	driverUC := driverUsecase.NewUseCase(driverRepo)
	userUC := userUsecase.NewUseCase(userRepo, driverUC)
	inventoryUC := inventoryUsecase.NewUseCase(inventoryRepo)
	orderUC := orderUsecase.NewUseCase(orderRepo, &inventoryadapter.UseCaseAdapter{UseCase: inventoryUC}, &useradapter.UseCaseAdapter{UseCase: userUC})
	deliveryUC := deliveryUsecase.NewUseCase(deliveryRepo, &orderadapter.UseCaseAdapter{UseCase: orderUC}, &driveradapter.UseCaseAdapter{UseCase: driverUC})

	// Combined cross-domain service
	orderService := application.NewOrderService(
		&useradapter.UseCaseAdapter{UseCase: userUC},
		&orderadapter.UseCaseAdapter{UseCase: orderUC},
		&driveradapter.UseCaseAdapter{UseCase: driverUC},
		&deliveryadapter.UseCaseAdapter{UseCase: deliveryUC},
		&inventoryadapter.UseCaseAdapter{UseCase: inventoryUC},
	)

	// Other usecases
	paymentUC := paymentUsecase.NewUseCase(paymentRepo)
	feedbackUC := feedbackUsecase.NewUseCase(feedbackRepo)
	notificationUC := notificationUsecase.NewUseCase(notificationRepo)

	// Set up Handlers
	userHandler := handlers.NewUserHandler(userUC)
	orderHandler := handlers.NewOrderHandler(orderService)
	driverHandler := handlers.NewDriverHandler(driverUC)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryUC)
	paymentHandler := handlers.NewPaymentHandler(paymentUC)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackUC)
	notificationHandler := handlers.NewNotificationHandler(notificationUC)
	inventoryHandler := handlers.NewInventoryHandler(inventoryUC)
	inviteHandler := handlers.NewInviteHandler(inviteUC)

	// Start server
	r := router.NewRouter(
		userHandler,
		orderHandler,
		driverHandler,
		deliveryHandler,
		paymentHandler,
		feedbackHandler,
		notificationHandler,
		inventoryHandler,
		publicApiBaseUrl,
		inviteHandler,
	)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("could not start server at: %v", err)
	}
}
