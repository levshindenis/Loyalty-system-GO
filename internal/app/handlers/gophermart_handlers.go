package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/storages"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/tools"
	"io"
	"net/http"
)

type HStorage struct {
	storages.ServerStorage
}

func (hs *HStorage) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var buf bytes.Buffer

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "Empty login or password", http.StatusBadRequest)
		return
	}

	flag, userId, err := hs.CheckUser(user.Login, user.Password, "registration")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if flag {
		http.Error(w, "Login used", http.StatusConflict)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "UserID",
		Value: userId,
	})
}

func (hs *HStorage) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var buf bytes.Buffer

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "Empty login or password", http.StatusBadRequest)
		return
	}

	flag, userId, err := hs.CheckUser(user.Login, user.Password, "login")
	if err != nil {
		http.Error(w, "Something bad with check user", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Bad pair login & password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "UserID",
		Value: userId,
	})
}

func (hs *HStorage) MakeOrderHandler(w http.ResponseWriter, r *http.Request) {
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

	flag, err := tools.IsLuna(string(body))
	if err != nil {
		http.Error(w, "Something bad with IsLuna", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Failed the Luna test", http.StatusUnprocessableEntity)
		return
	}

	cookie, _ := r.Cookie("UserID")

	fl1, fl2, err := hs.CheckOrder(string(body), cookie.Value, "make")
	if fl1 && !fl2 {
		http.Error(w, "Order made other person", http.StatusConflict)
		return
	}
	if !fl1 {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (hs *HStorage) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	isFilled, ordersData, err := hs.GetOrders(cookie.Value)
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

func (hs *HStorage) CountPointsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	balanceData, err := hs.GetBalance(cookie.Value)
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

func (hs *HStorage) DeductPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	var withdraw models.Withdraw
	var buf bytes.Buffer

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

	flag, _, err := hs.CheckOrder(withdraw.OrderId, cookie.Value, "check")
	if err != nil {
		http.Error(w, "Something bad with CheckOrder", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Bad order value", http.StatusUnprocessableEntity)
		return
	}

	flag, err = hs.CheckBalance(cookie.Value, withdraw.OrderId, withdraw.Summ)
	if err != nil {
		http.Error(w, "Something bad with CheckBalance", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Not enough bonuses", http.StatusPaymentRequired)
		return
	}
}

func (hs *HStorage) MovementPointsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	flag, outPointsData, err := hs.GetOutPoints(cookie.Value)
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
