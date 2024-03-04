package main

import (
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/routers"
)

func main() {
	var hs handlers.HStorage
	if err := initHS(&hs); err != nil {
		panic(err)
	}
	if err := http.ListenAndServe(hs.GetRunAddress(), routers.GophermartRouter(hs)); err != nil {
		panic(err)
	}
}

func initHS(hs *handlers.HStorage) error {
	hs.ParseFlags()
	if err := hs.Init(); err != nil {
		return err
	}
	return nil
}
