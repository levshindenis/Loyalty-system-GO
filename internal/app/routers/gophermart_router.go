package routers

import (
	"github.com/go-chi/chi/v5"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
)

func GophermartRouter(hs handlers.HStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hs.RegisterHandler)
		r.Post("/login", hs.LoginHandler)
		r.Post("/orders", hs.MakeOrderHandler)
		r.Get("/orders", hs.GetOrdersHandler)
		r.Get("/balance", hs.CountPointsHandler)
		r.Post("/balance/withdraw", hs.DeductPointsHandler)
		r.Get("/withdrawals", hs.MovementPointsHandler)
	})
	return r
}
