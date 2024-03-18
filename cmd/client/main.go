package main

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/looplab/fsm"
)

type Server struct {
	client  *resty.Client
	cookie  string
	address string
	choice  string
	m       map[string]string
	f       *fsm.FSM
}

func NewServer() *Server {
	client := resty.New()
	m := map[string]string{
		"1": "reg",
		"2": "log",
		"3": "mOrd",
		"4": "gOrd",
		"5": "gBal",
		"6": "setWd",
		"7": "getWd",
	}
	return &Server{
		client:  client,
		cookie:  "",
		address: "http://localhost:8000",
		choice:  "",
		m:       m,
	}
}

func main() {
	server := NewServer()
	server.f = fsm.NewFSM(
		"zero",
		fsm.Events{
			{Name: "go", Src: []string{"zero"}, Dst: "start"},
			{Name: "reg", Src: []string{"start"}, Dst: "register"},
			{Name: "log", Src: []string{"start"}, Dst: "login"},
			{Name: "mOrd", Src: []string{"main"}, Dst: "makeOrder"},
			{Name: "gOrd", Src: []string{"main"}, Dst: "getOrders"},
			{Name: "gBal", Src: []string{"main"}, Dst: "getBalance"},
			{Name: "setWd", Src: []string{"main"}, Dst: "setWithdraw"},
			{Name: "getWd", Src: []string{"main"}, Dst: "getWithdrawals"},
			{Name: "mainpage",
				Src: []string{"register", "login", "makeOrder", "getOrders", "getBalance", "setWithdraw", "getWithdrawals"},
				Dst: "main"},
		},
		fsm.Callbacks{
			"start":    func(_ context.Context, _ *fsm.Event) { server.RegLog() },
			"reg":      func(_ context.Context, _ *fsm.Event) { server.Auth("register") },
			"log":      func(_ context.Context, _ *fsm.Event) { server.Auth("login") },
			"mOrd":     func(_ context.Context, _ *fsm.Event) { server.MakeOrder() },
			"gOrd":     func(_ context.Context, _ *fsm.Event) { server.GetOrders() },
			"gBal":     func(_ context.Context, _ *fsm.Event) { server.GetBalance() },
			"setWd":    func(_ context.Context, _ *fsm.Event) { server.SetWithdraw() },
			"getWd":    func(_ context.Context, _ *fsm.Event) { server.GetWithdrawals() },
			"mainpage": func(_ context.Context, _ *fsm.Event) { server.SelectAction() },
		},
	)

	if err := server.f.Event(context.Background(), "go"); err != nil {
		panic(err)
	}
}
