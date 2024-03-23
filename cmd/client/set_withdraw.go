package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (s *Server) SetWithdraw() {
	var order string
	var summ float64

	fmt.Println("Введите номер заказа:")
	fmt.Scanf("%s\n", &order)
	fmt.Println("Введите сумму списания:")
	fmt.Scanf("%f\n", &summ)

	withDraw := models.Withdraw{OrderID: order, Summ: summ}

	jsonWD, err := json.Marshal(withDraw)
	if err != nil {
		panic(err)
	}

	resp, err := s.client.R().SetBody(bytes.NewBuffer(jsonWD)).
		SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).Post(s.address + "/api/user/balance/withdraw")
	if err != nil {
		panic(err)
	}

	fmt.Println("Ответ:")
	fmt.Println(resp.Status())

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
