package main

import (
	"context"
	"database/sql"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/config"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	conf config.ServerConfig
	hs   handlers.HStorage
)

func TestMain(m *testing.M) {
	conf.ParseFlags()
	db, err := sql.Open("pgx", conf.GetDBURI())
	if err != nil {
		log.Fatalf("Error with Open")
	}
	defer db.Close()

	if err = hs.Init(db, conf.GetAccSysAddr()); err != nil {
		log.Fatalf("Error with Init")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error with Begin")
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (user_id, login, password) values ($1, $2, $3)`,
		"abc", "aaa", "bbb")
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with INSERT INTO users")
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (user_id, login, password) values ($1, $2, $3)`,
		"def", "bbb", "ccc")
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with INSERT INTO users")
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO balances (user_id, balance, withdrawn) values ($1, $2, $3)`,
		"abc", 10, 5)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with INSERT INTO balances")
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO balances (user_id, balance, withdrawn) values ($1, $2, $3)`,
		"def", 0, 0)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with INSERT INTO balances")
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
		"1115", "abc", "PROCESSED", 15, time.Now().Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with INSERT INTO orders")
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
		"2220", "abc", "SOLD", 5, time.Now().Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with INSERT INTO orders")
	}
	tx.Commit()

	exitVal := m.Run()

	tx, err = db.Begin()
	_, err = tx.ExecContext(ctx,
		`DELETE from orders where user_id = $1`, "abc")
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with DELETE orders")
	}

	_, err = tx.ExecContext(ctx,
		`DELETE from users where login in ($1, $2, $3)`, "aaa", "123", "bbb")
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with DELETE orders")
	}

	_, err = tx.ExecContext(ctx,
		`DELETE from balances where user_id in($1, $2) or balance = 0`, "abc", "def")
	if err != nil {
		tx.Rollback()
		log.Fatal("Error with DELETE orders")
	}
	tx.Commit()

	os.Exit(exitVal)
}

func TestHStorage_Register(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		requestBody  string
		contentType  string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			requestBody:  "{\"login\":\"123\",\"password\":\"123\"}",
			contentType:  "application/json",
			expectedCode: 200,
		},
		{
			name:         "Repeat login",
			method:       http.MethodPost,
			requestBody:  "{\"login\":\"aaa\",\"password\":\"bbc\"}",
			contentType:  "application/json",
			expectedCode: 409,
		},
		{
			name:         "Bad request",
			method:       http.MethodPost,
			requestBody:  "{\"logn\":\"aaa\",\"password\":\"bbc\"}",
			contentType:  "application/json",
			expectedCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/register", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", tt.contentType)
			hs.Register(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHStorage_Login(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		requestBody  string
		contentType  string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			requestBody:  "{\"login\":\"aaa\",\"password\":\"bbb\"}",
			contentType:  "application/json",
			expectedCode: 200,
		},
		{
			name:         "Bad pair",
			method:       http.MethodPost,
			requestBody:  "{\"login\":\"aza\",\"password\":\"bbc\"}",
			contentType:  "application/json",
			expectedCode: 401,
		},
		{
			name:         "Bad request",
			method:       http.MethodPost,
			requestBody:  "{\"logn\":\"aaa\",\"password\":\"bbc\"}",
			contentType:  "application/json",
			expectedCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/login", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", tt.contentType)
			hs.Login(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHStorage_MakeOrder(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		requestBody  string
		cookie       string
		contentType  string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			requestBody:  "3335",
			cookie:       "abc",
			contentType:  "text/plain",
			expectedCode: 202,
		},
		{
			name:         "Self order repeat",
			method:       http.MethodPost,
			requestBody:  "1115",
			cookie:       "abc",
			contentType:  "text/plain",
			expectedCode: 200,
		},
		{
			name:         "Another order repeat",
			method:       http.MethodPost,
			requestBody:  "1115",
			cookie:       "def",
			contentType:  "text/plain",
			expectedCode: 409,
		},
		{
			name:         "Bad Luna test",
			method:       http.MethodPost,
			requestBody:  "3331",
			contentType:  "text/plain",
			expectedCode: 422,
		},
		{
			name:         "Bad request",
			method:       http.MethodPost,
			requestBody:  "{\"logn\":\"aaa\",\"password\":\"bbc\"}",
			contentType:  "application/json",
			expectedCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/orders", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", tt.contentType)
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookie})
			hs.MakeOrder(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHStorage_GetOrders(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		cookie       string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodGet,
			cookie:       "abc",
			expectedCode: 200,
		},
		{
			name:         "No orders",
			method:       http.MethodGet,
			cookie:       "def",
			expectedCode: 204,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/orders", nil)
			w := httptest.NewRecorder()
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookie})
			hs.GetOrders(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHStorage_CountPoints(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		cookie       string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodGet,
			cookie:       "abc",
			expectedCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/balance", nil)
			w := httptest.NewRecorder()
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookie})
			hs.CountPoints(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHStorage_DeductPoints(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		requestBody  string
		cookie       string
		contentType  string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			requestBody:  "{\"order\":\"4440\",\"sum\":5}",
			cookie:       "abc",
			contentType:  "application/json",
			expectedCode: 200,
		},
		{
			name:         "No money",
			method:       http.MethodPost,
			requestBody:  "{\"order\":\"5553\",\"sum\":100}",
			cookie:       "abc",
			contentType:  "application/json",
			expectedCode: 402,
		},
		{
			name:         "Bad order number",
			method:       http.MethodPost,
			requestBody:  "{\"order\":\"1113\",\"sum\":100}",
			cookie:       "abc",
			contentType:  "application/json",
			expectedCode: 422,
		},
		{
			name:         "Server error",
			method:       http.MethodPost,
			requestBody:  "1113",
			cookie:       "abc",
			contentType:  "application/json",
			expectedCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/balance/withdraw", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", tt.contentType)
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookie})
			hs.DeductPoints(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHStorage_MovementPoints(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		cookie       string
		expectedCode int
	}{
		{
			name:         "Good test",
			method:       http.MethodGet,
			cookie:       "abc",
			expectedCode: 200,
		},
		{
			name:         "Good test",
			method:       http.MethodGet,
			cookie:       "def",
			expectedCode: 204,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/user/withdrawals", nil)
			w := httptest.NewRecorder()
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookie})
			hs.MovementPoints(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}
