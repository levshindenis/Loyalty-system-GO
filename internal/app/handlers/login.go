package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (hs *HStorage) Login(w http.ResponseWriter, r *http.Request) {
	var (
		user models.User
		buf  bytes.Buffer
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

	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "Empty login or password", http.StatusBadRequest)
		return
	}

	flag, userID, err := hs.CheckUser(user.Login, user.Password, "login")
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
		Value: userID,
	})
}
