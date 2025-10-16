package router

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"logistics-backend/handlers"
	authMiddleware "logistics-backend/internal/middleware"
)

func NewRouter(
	u *handlers.UserHandler,
	o *handlers.OrderHandler,
	d *handlers.DriverHandler,
	e *handlers.DeliveryHandler,
	p *handlers.PaymentHandler, f *handlers.FeedbackHandler,
	n *handlers.NotificationHandler,
	i *handlers.InventoryHandler,
	publicApiBaseUrl string,
	c *handlers.InviteHandler,
	s *handlers.StoreHandler,
) http.Handler {
	r := chi.NewRouter()

	// Enable Cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		// Swagger docs (this will now be served under /api/swagger)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(publicApiBaseUrl+"/swagger/doc.json"),
		))

		// Public routes
		r.Route("/public", func(r chi.Router) {
			// Public auth
			r.Post("/create", u.CreateUser)
			r.Post("/login", u.LoginUser)

			// Public store pages
			r.Route("/stores", func(r chi.Router) {
				r.Get("/public", s.GetPublicStores)
			})
		})

		// Protected Routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.JWTAuthMiddleware)

			// Users
			r.Route("/users", func(r chi.Router) {
				r.Get("/all_users", u.ListUsers)
				r.Get("/by-id/{id}", u.GetUserByID)
				r.Get("/by-email/{email}", u.GetUserByEmail)
				r.Patch("/{id}/profile", u.UpdateUserProfile)
				r.Put("/{id}/update", u.UpdateUser)
				r.Delete("/{id}", u.DeleteUser)
			})

			// Invites
			r.Route("/invites", func(r chi.Router) {
				r.Post("/create", c.CreateMember)
				r.Get("/by-token", c.GetMemberByToken)
				r.Get("/all_invites", c.ListPendingMembers)
				r.Delete("/{id}", c.DeleteMember)
			})

			// Orders
			r.Route("/orders", func(r chi.Router) {
				r.Post("/create", o.CreateOrder)
				r.Get("/all_orders", o.ListOrders)
				r.Get("/form-data", o.GetOrderFormData)
				r.Post("/assign", o.AutoAssignOrders)
				r.Get("/by-id/{id}", o.GetOrderByID)
				r.Get("/by-customer/{customer_id}", o.GetOrderByCustomer)
				r.Put("/{id}/update", o.UpdateOrder)
				r.Delete("/{id}", o.DeleteOrder)
			})

			// Inventories
			r.Route("/inventories", func(r chi.Router) {
				r.Post("/create", i.CreateInventory)
				r.Get("/by-name", i.GetByInventoryName)
				r.Get("/all_inventories", i.ListInventories)
				r.Get("/by-category", i.GetInventoryByCategory)
				r.Get("/categories", i.ListCategories)
				r.Get("/by-id/{id}", i.GetByInventoryID)
				r.Get("/by-store/{id}", i.GetInventoryByStore)
				r.Delete("/{id}", i.DeleteInventory)
			})

			// Drivers
			r.Route("/drivers", func(r chi.Router) {
				r.Get("/all_drivers", d.ListDrivers)
				r.Get("/by-id/{id}", d.GetDriverByID)
				r.Get("/by-email/{email}", d.GetDriverByEmail)
				r.Patch("/{id}/profile", d.UpdateDriverProfile)
				r.Put("/{id}/update", d.UpdateDriver)
				r.Delete("/{id}", d.DeleteDriver)
			})

			// Deliveries
			r.Route("/deliveries", func(r chi.Router) {
				r.Get("/all_deliveries", e.ListDeliveries)
				r.Get("/by-id/{id}", e.GetDeliveryByID)
				r.Put("/{id}/update", e.UpdateDelivery)
				r.Put("/{id}/accept", e.AcceptDelivery)
				r.Delete("/{id}", e.DeleteDelivery)
			})

			// Payments
			r.Route("/payments", func(r chi.Router) {
				r.Post("/create", p.CreatePayment)
				r.Get("/all_payments", p.ListPayments)
				r.Get("/{id}", p.GetPaymentByID)
				r.Get("/{order_id}", p.GetPaymentByOrderID)
			})

			// Feedbacks
			r.Route("/feedbacks", func(r chi.Router) {
				r.Post("/create", f.CreateFeedback)
				r.Get("/all_feedbacks", f.ListFeedback)
				r.Get("/{id}", f.GetFeedbackByID)
			})

			// Notifications
			r.Route("/notifications", func(r chi.Router) {
				r.Post("/create", n.CreateNotification)
				r.Get("/all_pending_notifications", n.ListNotifications)
				r.Get("/all_my_notifications/{id}", n.ListUserNotifications)
				r.Put("/{id}/status", n.UpdateNotificationStatus)
				r.Patch("/{id}/read", n.MarkAsRead)
				r.Patch("/mark_all_as_read/{id}", n.MarkAllAsRead)
			})

			// Stores
			r.Route("/stores", func(r chi.Router) {
				r.Post("/create", s.CreateStore)
				r.Get("/by-slug", s.GetStoreBySlug)
				r.Get("/public", s.GetPublicStores)
				r.Get("/by-id/{id}", s.GetStoreByID)
				r.Put("/{id}/update", s.UpdateStore)
				r.Delete("/{id}", s.DeleteStore)
			})

		})
	})

	return r
}
