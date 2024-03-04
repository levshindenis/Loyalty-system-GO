package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	cookVal := ""

	for {
		var choice string
		client := resty.New()

		SelectAction()

		fmt.Scanf("%s\n", &choice)

		fmt.Println()

		switch choice {
		case "1":
			value, err := Register(client)
			if err != nil {
				panic(err)
			}
			cookVal = value
		case "2":
			value, err := Login(client)
			if err != nil {
				panic(err)
			}
			cookVal = value
		case "3":
			err := MakeOrder(client, cookVal)
			if err != nil {
				panic(err)
			}
		case "4":
			err := GetOrders(client, cookVal)
			if err != nil {
				panic(err)
			}
		case "5":
			err := GetBalance(client, cookVal)
			if err != nil {
				panic(err)
			}
		case "6":
			err := SetWithdraw(client, cookVal)
			if err != nil {
				panic(err)
			}
		case "7":
			err := GetWithdrawals(client, cookVal)
			if err != nil {
				panic(err)
			}
		case "8":
			err := SetOrder(client)
			if err != nil {
				panic(err)
			}
		default:
			os.Exit(1)
		}

		Exit()

		fmt.Scanf("%s\n", &choice)

		switch choice {
		case "1":
			os.Exit(0)
		default:
			fmt.Println()
			continue
		}
	}
	fmt.Println(cookVal)
}

func SelectAction() {
	fmt.Println("Выберите действие:")
	fmt.Println("1) Регистрация")
	fmt.Println("2) Вход")
	fmt.Println("3) Создать заказ")
	fmt.Println("4) Получить список заказов")
	fmt.Println("5) Баланс")
	fmt.Println("6) Списать баллы")
	fmt.Println("7) Выводы баллов")
	fmt.Println("8) Добавить заказ в систему рассчета")
	fmt.Println("===========================")
	fmt.Print("Ввод:  ")
}

func Exit() {
	fmt.Println("\nВыход:")
	fmt.Println("1) Да")
	fmt.Println("2) Нет")
	fmt.Println("=============")
	fmt.Print("Ввод:  ")
}

func Register(client *resty.Client) (string, error) {
	var login, password string
	fmt.Println("Введите логин:   ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Введите пароль:  ")
	fmt.Scanf("%s\n", &password)

	user := models.User{Login: login, Password: password}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	resp, err := client.R().SetBody(bytes.NewBuffer(jsonUser)).Post("http://localhost:8000/api/user/register")
	if err != nil {
		return "", err
	}

	fmt.Println("Ответ:")
	if resp.StatusCode() == 200 {
		fmt.Println("All good!")
		return resp.Cookies()[0].Value, nil
	}

	fmt.Println(resp.Status())
	return "", nil
}

func Login(client *resty.Client) (string, error) {
	var login, password string
	fmt.Println("Введите логин:   ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Введите пароль:  ")
	fmt.Scanf("%s\n", &password)

	user := models.User{Login: login, Password: password}

	jsonUser, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		return "", err
	}

	resp, err := client.R().SetBody(bytes.NewBuffer(jsonUser)).Post("http://localhost:8000/api/user/login")
	if err != nil {
		return "", err
	}

	fmt.Println("Ответ:")
	if resp.StatusCode() == 200 {
		fmt.Println("All good!")
		return resp.Cookies()[0].Value, nil
	}

	fmt.Println(resp.Status())
	return "", nil

}

func MakeOrder(client *resty.Client, cookie string) error {
	var order string
	fmt.Println("Введите номер заказа:")
	fmt.Scanf("%s\n", &order)

	resp, err := client.R().SetHeader("Content-Type", "text/plain").SetBody(order).SetCookie(&http.Cookie{Name: "UserID", Value: cookie}).Post("http://localhost:8000/api/user/orders")
	if err != nil {
		return err
	}

	fmt.Println("Ответ:")
	if resp.StatusCode() == 202 {
		fmt.Println("All good!")
	} else {
		fmt.Println(resp.Status())
	}
	return nil
}

func GetOrders(client *resty.Client, cookie string) error {
	fmt.Println("Ответ:")

	resp, err := client.R().SetCookie(&http.Cookie{Name: "UserID", Value: cookie}).Get("http://localhost:8000/api/user/orders")
	if err != nil {
		return err
	}

	if resp.StatusCode() == 200 {
		fmt.Println(resp)
	} else {
		fmt.Println(resp.Status())
	}

	return nil
}

func GetBalance(client *resty.Client, cookie string) error {
	fmt.Println("Ответ:")

	resp, err := client.R().SetCookie(&http.Cookie{Name: "UserID", Value: cookie}).Get("http://localhost:8000/api/user/balance")
	if err != nil {
		return err
	}

	if resp.StatusCode() == 200 {
		fmt.Println(resp)
	} else {
		fmt.Println(resp.Status())
	}

	return nil
}

func SetWithdraw(client *resty.Client, cookie string) error {
	var order string
	var summ float64

	fmt.Println("Введите номер заказа:")
	fmt.Scanf("%s\n", &order)
	fmt.Println("Введите сумму списания:")
	fmt.Scanf("%f\n", &summ)

	withDraw := models.Withdraw{OrderID: order, Summ: summ}

	jsonWD, err := json.Marshal(withDraw)
	if err != nil {
		return err
	}

	resp, err := client.R().SetBody(bytes.NewBuffer(jsonWD)).SetCookie(&http.Cookie{Name: "UserID", Value: cookie}).Post("http://localhost:8000/api/user/balance/withdraw")
	if err != nil {
		return err
	}

	fmt.Println("Ответ:")
	if resp.StatusCode() == 200 {
		fmt.Println("All Good!")
	} else {
		fmt.Println(resp.Status())
	}

	return nil
}

func GetWithdrawals(client *resty.Client, cookie string) error {
	fmt.Println("Ответ:")

	resp, err := client.R().SetCookie(&http.Cookie{Name: "UserID", Value: cookie}).Get("http://localhost:8000/api/user/withdrawals")
	if err != nil {
		return err
	}

	if resp.StatusCode() == 200 {
		fmt.Println(resp)
	} else {
		fmt.Println(resp.Status())
	}

	return nil
}

//

func SetOrder(client *resty.Client) error {
	var orderID, name string
	var prices []float64
	var percentages []float64

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
