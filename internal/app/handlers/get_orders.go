package handlers

import (
	"encoding/json"
	"net/http"
)

func (hs *HStorage) GetOrders(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	isFilled, ordersData, err := hs.GetUserOrders(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isFilled {
		http.Error(w, "Empty orders data", http.StatusNoContent)
		return
	}

	resp, err := json.MarshalIndent(ordersData, "", "    ")
	if err != nil {
		http.Error(w, "Something bad with Marshal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(resp); err != nil {
		http.Error(w, "Something bad with write to ResponseWriter", http.StatusInternalServerError)
		return
	}
}
