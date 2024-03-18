package handlers

import (
	"encoding/json"
	"net/http"
)

func (hs *HStorage) MovementPoints(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	flag, outPointsData, err := hs.GetUserOutPoints(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Not enough orders", http.StatusNoContent)
		return
	}

	resp, err := json.MarshalIndent(outPointsData, "", "    ")
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
