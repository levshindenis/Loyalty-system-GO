package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/luna"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (hs *HStorage) DeductPoints(w http.ResponseWriter, r *http.Request) {
	var (
		withdraw models.Withdraw
		buf      bytes.Buffer
	)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &withdraw); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusInternalServerError)
		return
	}

	cookie, _ := r.Cookie("UserID")

	flag, err := luna.IsLuna(withdraw.OrderID)
	if err != nil {
		http.Error(w, "Something bad with CheckUserOrder", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Bad order value", http.StatusUnprocessableEntity)
		return
	}

	flag, err = hs.CheckUserBalance(cookie.Value, withdraw.OrderID, withdraw.Summ)
	if err != nil {
		http.Error(w, "Something bad with CheckUserBalance", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Not enough bonuses", http.StatusPaymentRequired)
		return
	}
}
