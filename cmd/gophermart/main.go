package main

import (
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/routers"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var hs handlers.HStorage
	hs.ParseFlags()
	hs.Init()

	return http.ListenAndServe(hs.GetRunAddress(), routers.GophermartRouter(hs))
}
