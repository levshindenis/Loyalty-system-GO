package main

import (
	"context"
	"fmt"
)

func (s *Server) RegLog() {
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1) Регистрация")
		fmt.Println("2) Вход")
		fmt.Println("==================")
		fmt.Print("Ввод:  ")
		fmt.Scanf("%s/n", &s.choice)
		if s.choice == "1" || s.choice == "2" {
			break
		}
		fmt.Println("Bad answer. Please repeat!")
	}

	if err := s.f.Event(context.Background(), s.m[s.choice]); err != nil {
		panic(err)
	}
}
