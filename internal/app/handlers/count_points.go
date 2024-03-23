package handlers

import (
	"encoding/json"
	"net/http"
)

func (hs *HStorage) CountPoints(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	balanceData, err := hs.GetUserBalance(cookie.Value)
	if err != nil {
		http.Error(w, "Something bad with check cookie", http.StatusInternalServerError)
		return
	}

	resp, err := json.MarshalIndent(balanceData, "", "    ")
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
