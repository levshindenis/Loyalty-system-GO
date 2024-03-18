package main

import (
	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()

	if err := SetOrder(client); err != nil {
		panic(err)
	}
}
