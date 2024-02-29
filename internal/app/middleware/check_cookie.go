package middleware

import (
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
)

func CheckCookie(next http.HandlerFunc, hs handlers.HStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("UserID")
		if err != nil {
			http.Error(w, "Not cookie", http.StatusUnauthorized)
			return
		}

		flag, err := hs.CheckCookie(cookie.Value)
		if err != nil {
			http.Error(w, "Something bad with check cookie", http.StatusInternalServerError)
			return
		}
		if !flag {
			http.Error(w, "Bad cookie", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
