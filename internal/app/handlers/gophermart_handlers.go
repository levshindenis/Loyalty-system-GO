package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type Decoder struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var dec Decoder
	var buf bytes.Buffer

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &dec); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusBadRequest)
		return
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type Decoder struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var dec Decoder
	var buf bytes.Buffer

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &dec); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusBadRequest)
		return
	}
}

func MakeOrderHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func CountPointsHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func DeductPointsHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func MovementPointsHandler(w http.ResponseWriter, r *http.Request) {
	//
}
