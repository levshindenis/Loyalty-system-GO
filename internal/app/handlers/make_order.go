package handlers

import (
	"github.com/levshindenis/Loyalty-system-GO/internal/app/luna"
	"io"
	"net/http"
)

func (hs *HStorage) MakeOrder(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something bad with read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	flag, err := luna.IsLuna(string(body))
	if err != nil {
		http.Error(w, "Something bad with IsLuna", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Failed the Luna test", http.StatusUnprocessableEntity)
		return
	}

	cookie, _ := r.Cookie("UserID")

	fl1, fl2, err := hs.CheckUserOrder(string(body), cookie.Value)
	if err != nil {
		http.Error(w, "Something bad with CheckUserOrder", http.StatusInternalServerError)
		return
	}
	if fl1 && !fl2 {
		http.Error(w, "Order made other person", http.StatusConflict)
		return
	}
	if !fl1 {
		w.WriteHeader(http.StatusAccepted)
	}
}
