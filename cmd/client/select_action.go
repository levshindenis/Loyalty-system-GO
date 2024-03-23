package main

import (
	"context"
	"fmt"
	"strconv"
)

func (s *Server) SelectAction() {
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1) Создать заказ")
		fmt.Println("2) Получить список заказов")
		fmt.Println("3) Баланс")
		fmt.Println("4) Списать баллы")
		fmt.Println("5) Выводы баллов")
		fmt.Println("===========================")
		fmt.Print("Ввод:  ")
		fmt.Scanf("%s/n", &s.choice)
		if s.choice == "1" || s.choice == "2" || s.choice == "3" || s.choice == "4" || s.choice == "5" {
			break
		}
		fmt.Println("Bad answer. Please repeat!")
	}

	ch, err := strconv.Atoi(s.choice)
	if err != nil {
		panic(err)
	}
	ch += 2

	if err = s.f.Event(context.Background(), s.m[strconv.Itoa(ch)]); err != nil {
		panic(err)
	}
}
