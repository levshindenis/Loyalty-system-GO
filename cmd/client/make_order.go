package main

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Server) MakeOrder() {
	var order string
	fmt.Println("Введите номер заказа:")
	fmt.Scanf("%s\n", &order)

	resp, err := s.client.R().SetHeader("Content-Type", "text/plain").
		SetBody(order).SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).Post(s.address + "/api/user/orders")
	if err != nil {
		panic(err)
	}

	fmt.Println("Ответ:")

	fmt.Println(resp.Status())

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
