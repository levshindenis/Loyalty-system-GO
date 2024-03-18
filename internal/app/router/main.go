package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/middleware"
)

func Router(hs handlers.HStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", hs.Register)
		r.Post("/login", hs.Login)
		r.Post("/orders", middleware.CheckCookie(hs.MakeOrder, hs))
		r.Get("/orders", middleware.CheckCookie(hs.GetOrders, hs))
		r.Get("/balance", middleware.CheckCookie(hs.CountPoints, hs))
		r.Post("/balance/withdraw", middleware.CheckCookie(hs.DeductPoints, hs))
		r.Get("/withdrawals", middleware.CheckCookie(hs.MovementPoints, hs))
	})
	return r
}
