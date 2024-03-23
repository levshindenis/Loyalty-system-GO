package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (s *Server) Auth(param string) {
	var login, password string
	fmt.Println("Введите логин:   ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Введите пароль:  ")
	fmt.Scanf("%s\n", &password)

	user := models.User{Login: login, Password: password}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	resp, err := s.client.R().SetBody(bytes.NewBuffer(jsonUser)).Post(s.address + "/api/user/" + param)
	if err != nil {
		panic(err)
	}

	s.cookie = resp.Cookies()[0].Value

	fmt.Println("Ответ:")

	fmt.Println(resp.Status())

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
