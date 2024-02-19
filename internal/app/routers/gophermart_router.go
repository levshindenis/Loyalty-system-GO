package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
)

func MyRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.RegisterHandler)
		r.Post("/login", handlers.LoginHandler)
		r.Post("/orders", handlers.MakeOrderHandler)
		r.Get("/orders", handlers.GetOrdersHandler)
		r.Get("/balance", handlers.CountPointsHandler)
		r.Post("/balance/withdraw", handlers.DeductPointsHandler)
		r.Get("/withdrawals", handlers.MovementPointsHandler)
	})
	return r
}
