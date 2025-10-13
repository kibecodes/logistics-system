package main

import (
	"log"
	"net/http"
	"os"

	"logistics-backend/handlers"
	deliveryadapter "logistics-backend/internal/adapters/delivery"
	driveradapter "logistics-backend/internal/adapters/driver"
	inventoryadapter "logistics-backend/internal/adapters/inventory"
	notificationadapter "logistics-backend/internal/adapters/notification"
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
	storeUsecase "logistics-backend/internal/usecase/store"
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

	db := sqlx.MustConnect("postgres", dbUrl)

	txm := application.NewTxManager(db)

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
	storeRepo := postgres.NewStoreRepository(db)

	// Set up usecase
	// Individual
	inviteUC := inviteUsecase.NewUseCase(inviteRepo, txm)
	driverUC := driverUsecase.NewUseCase(driverRepo, txm, notificationRepo)
	userUC := userUsecase.NewUseCase(userRepo, driverUC, txm, notificationRepo)
	inventoryUC := inventoryUsecase.NewUseCase(inventoryRepo, txm, notificationRepo, storeRepo)
	orderUC := orderUsecase.NewUseCase(orderRepo, &inventoryadapter.UseCaseAdapter{UseCase: inventoryUC}, &useradapter.UseCaseAdapter{UseCase: userUC}, txm, notificationRepo, storeRepo)
	deliveryUC := deliveryUsecase.NewUseCase(deliveryRepo, &orderadapter.UseCaseAdapter{UseCase: orderUC}, &driveradapter.UseCaseAdapter{UseCase: driverUC}, txm, notificationRepo)
	notificationUC := notificationUsecase.NewUseCase(notificationRepo, txm)
	storeUC := storeUsecase.NewUseCase(storeRepo, txm)

	// Combined cross-domain service
	orderService := application.NewOrderService(
		&useradapter.UseCaseAdapter{UseCase: userUC},
		&orderadapter.UseCaseAdapter{UseCase: orderUC},
		&driveradapter.UseCaseAdapter{UseCase: driverUC},
		&deliveryadapter.UseCaseAdapter{UseCase: deliveryUC},
		&inventoryadapter.UseCaseAdapter{UseCase: inventoryUC},
		&notificationadapter.UseCaseAdapter{UseCase: notificationUC},
	)

	// Other usecases
	paymentUC := paymentUsecase.NewUseCase(paymentRepo, txm)
	feedbackUC := feedbackUsecase.NewUseCase(feedbackRepo, txm)

	// Set up Handlers
	userHandler := handlers.NewUserHandler(orderService)
	orderHandler := handlers.NewOrderHandler(orderService)
	driverHandler := handlers.NewDriverHandler(orderService)
	deliveryHandler := handlers.NewDeliveryHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentUC)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackUC)
	notificationHandler := handlers.NewNotificationHandler(orderService)
	inviteHandler := handlers.NewInviteHandler(inviteUC)
	inventoryHandler := handlers.NewInventoryHandler(orderService)
	storeHandler := handlers.NewStoreHandler(storeUC)

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
		storeHandler,
	)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("could not start server at: %v", err)
	}
}
