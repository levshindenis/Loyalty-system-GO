package main

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Server) GetBalance() {
	fmt.Println("Ответ:")

	resp, err := s.client.R().SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).Get(s.address + "/api/user/balance")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
