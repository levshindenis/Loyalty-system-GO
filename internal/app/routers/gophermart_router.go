package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/middleware"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
)

func GophermartRouter(hs handlers.HStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hs.RegisterHandler)
		r.Post("/login", hs.LoginHandler)
		r.Post("/orders", middleware.CheckCookie(hs.MakeOrderHandler, hs))
		r.Get("/orders", middleware.CheckCookie(hs.GetOrdersHandler, hs))
		r.Get("/balance", middleware.CheckCookie(hs.CountPointsHandler, hs))
		r.Post("/balance/withdraw", middleware.CheckCookie(hs.DeductPointsHandler, hs))
		r.Get("/withdrawals", middleware.CheckCookie(hs.MovementPointsHandler, hs))
	})
	return r
}
