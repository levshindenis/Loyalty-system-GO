package main

import (
	"github.com/levshindenis/Loyalty-system-GO/internal/app/config"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/routers"
	"net/http"
)

func main() {
	var perem config.ServerConfig
	perem.ParseFlags()

	if err := http.ListenAndServe("localhost:8080", routers.MyRouter()); err != nil {
		panic(err)
	}
}
