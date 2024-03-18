package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func SetOrder(client *resty.Client) error {
	var (
		orderID, name string
		prices        []float64
		percentages   []float64
	)

	for i := 1; i <= 5; i++ {
		prices = append(prices, float64(i*11))
	}

	for i := 1; i <= 5; i++ {
		percentages = append(percentages, float64(i*10))
	}

	fmt.Println("Номер заказа:")
	fmt.Scanf("%s\n", &orderID)
	fmt.Println("Название товара:")
	fmt.Scanf("%s\n", &name)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	randPerc := percentages[r.Intn(len(percentages))]
	fmt.Println("Процент:  ", randPerc)
	good := models.Good{Match: name, Reward: randPerc, RewardType: "%"}

	s = rand.NewSource(time.Now().Unix())
	r = rand.New(s)

	randPrice := prices[r.Intn(len(prices))]
	fmt.Println("Сумма:  ", randPrice)
	product := models.Product{Description: name, Price: randPrice}
	products := []models.Product{product}

	regOrder := models.RegOrder{OrderID: orderID, Products: products}

	jsonGood, err := json.Marshal(good)
	if err != nil {
		return err
	}

	resp1, err := client.R().SetBody(bytes.NewBuffer(jsonGood)).Post("http://localhost:8080/api/goods")
	if err != nil {
		return err
	}

	fmt.Println("Ответ на Goods:  ", resp1.StatusCode())

	jsonRegOrder, err := json.Marshal(regOrder)
	if err != nil {
		return err
	}

	resp2, err := client.R().SetBody(bytes.NewBuffer(jsonRegOrder)).Post("http://localhost:8080/api/orders")
	if err != nil {
		return err
	}

	fmt.Println("Ответ на Orders:  ", resp2.StatusCode())

	return nil
}
