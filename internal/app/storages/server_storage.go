package storages

import (
	"fmt"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/config"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
	"sync"
	"time"
)

type ServerStorage struct {
	sc     config.ServerConfig
	ds     DBStorage
	fromDB models.Queue
	toDB   models.Queue
}

func (serv *ServerStorage) ParseFlags() {
	serv.sc.ParseFlags()
}

func (serv *ServerStorage) Init() error {
	serv.ds = DBStorage{address: serv.sc.GetDBURI()}

	if err := serv.ds.MakeDB(); err != nil {
		return err
	}

	serv.fromDB = models.NewQueue()
	serv.toDB = models.NewQueue()

	go serv.FromDBToChannel(&serv.fromDB)

	m := sync.Mutex{}
	c := sync.NewCond(&m)
	for i := 0; i < 5; i++ {
		w := models.NewCompareWorker(i, serv.fromDB, serv.toDB, serv.sc.GetAccSysAddr(), c)
		go w.Loop()
	}

	go serv.FromChannelToDB(&serv.toDB)

	return nil
}

func (serv *ServerStorage) GetRunAddress() string {
	return serv.sc.GetRunAddress()
}

func (serv *ServerStorage) SetDBAddress(value string) {
	serv.ds.SetAddress(value)
}

func (serv *ServerStorage) CheckUser(login string, password string, param string) (bool, string, error) {
	return serv.ds.CheckUser(login, password, param)
}

func (serv *ServerStorage) CheckCookie(cookie string) (bool, error) {
	return serv.ds.CheckCookie(cookie)
}

func (serv *ServerStorage) CheckOrder(orderID string, userID string) (bool, bool, error) {
	return serv.ds.CheckOrder(orderID, userID)
}

func (serv *ServerStorage) GetOrders(userID string) (bool, []models.Order, error) {
	return serv.ds.GetOrders(userID)
}

func (serv *ServerStorage) GetBalance(userID string) (models.Balance, error) {
	return serv.ds.GetBalance(userID)
}

func (serv *ServerStorage) CheckBalance(userID string, orderID string, orderSum float64) (bool, error) {
	return serv.ds.CheckBalance(userID, orderID, orderSum)
}

func (serv *ServerStorage) GetOutPoints(userID string) (bool, []models.OutPoints, error) {
	return serv.ds.GetOutPoints(userID)
}

func (serv *ServerStorage) FromDBToChannel(q *models.Queue) {
	ticker := time.NewTicker(4 * time.Second)

	for {
		select {
		case <-q.GetCtx().Done():
			q.Push(models.Task{})
			return
		case <-ticker.C:
			items, err := serv.ds.GetNewOrders()
			if err != nil {
				panic(err)
			}
			for _, elem := range items {
				q.Push(elem)
			}
		}
	}
}

func (serv *ServerStorage) FromChannelToDB(q *models.Queue) {
	ticker := time.NewTicker(2 * time.Second)

	var values []models.Task

	for {
		select {
		case <-q.GetCtx().Done():
			q.Push(models.Task{})
			return
		case value := <-q.GetChannel():
			fmt.Println("Value: ", value)
			values = append(values, value)
		case <-ticker.C:
			if len(values) == 0 {
				continue
			}
			err := serv.ds.UpdateOrders(values)
			if err != nil {
				panic(err)
			}
			values = nil
		}
	}
}
