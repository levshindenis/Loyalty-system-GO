package main

import (
	"net/http"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/config"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/handlers"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/router"
)

func main() {
	var (
		hs handlers.HStorage
		sc config.ServerConfig
	)

	sc.ParseFlags()
	if err := hs.Init(sc.GetDBURI(), sc.GetAccSysAddr()); err != nil {
		panic(err)
	}
	if err := http.ListenAndServe(sc.GetRunAddress(), router.Router(hs)); err != nil {
		panic(err)
	}
	hs.Terminate()
}
