package storages

import (
	"fmt"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/accrual"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/config"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

type ServerStorage struct {
	sc     config.ServerConfig
	sd     ServerData
	sl     zap.SugaredLogger
	fromDB models.Queue
	toDB   models.Queue
}

func (serv *ServerStorage) ParseFlags() {
	serv.sc.ParseFlags()
}

func (serv *ServerStorage) InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	serv.sl = *logger.Sugar()
}

func (serv *ServerStorage) Init() error {
	db := DBStorage{address: serv.sc.GetDBURI()}

	if err := db.MakeDB(); err != nil {
		return err
	}

	serv.sd = ServerData{data: &db}

	serv.InitLogger()

	serv.fromDB = models.NewQueue()
	serv.toDB = models.NewQueue()

	go serv.FromDBToChannel(&serv.fromDB)

	m := sync.Mutex{}
	c := sync.NewCond(&m)
	for i := 0; i < 5; i++ {
		w := accrual.NewCompareWorker(i, serv.fromDB, serv.toDB, serv.sc.GetAccSysAddr(), c)
		go w.Loop(&serv.sl)
	}

	go serv.FromChannelToDB(&serv.toDB)

	return nil
}

func (serv *ServerStorage) Terminate() {
	serv.fromDB.GetCtx().Done()
	serv.toDB.GetCtx().Done()
}

func (serv *ServerStorage) GetRunAddress() string {
	return serv.sc.GetRunAddress()
}

func (serv *ServerStorage) CheckUser(login string, password string, param string) (bool, string, error) {
	return serv.sd.data.CheckUser(login, password, param)
}

func (serv *ServerStorage) CheckCookie(cookie string) (bool, error) {
	return serv.sd.data.CheckCookie(cookie)
}

func (serv *ServerStorage) CheckOrder(orderID string, userID string) (bool, bool, error) {
	return serv.sd.data.CheckOrder(orderID, userID)
}

func (serv *ServerStorage) GetOrders(userID string) (bool, []models.Order, error) {
	return serv.sd.data.GetOrders(userID)
}

func (serv *ServerStorage) GetBalance(userID string) (models.Balance, error) {
	return serv.sd.data.GetBalance(userID)
}

func (serv *ServerStorage) CheckBalance(userID string, orderID string, orderSum float64) (bool, error) {
	return serv.sd.data.CheckBalance(userID, orderID, orderSum)
}

func (serv *ServerStorage) GetOutPoints(userID string) (bool, []models.OutPoints, error) {
	return serv.sd.data.GetOutPoints(userID)
}

func (serv *ServerStorage) FromDBToChannel(q *models.Queue) {
	ticker := time.NewTicker(4 * time.Second)

	for {
		select {
		case <-q.GetCtx().Done():
			q.Push(models.Task{})
			return
		case <-ticker.C:
			items, err := serv.sd.data.GetNewOrders()
			if err != nil {
				serv.sl.Infoln(
					"time", time.Now(),
					"error", "Error with FromDBChannel",
				)
				continue
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
			err := serv.sd.data.UpdateOrders(values)
			if err != nil {
				serv.sl.Infoln(
					"time", time.Now(),
					"error", "Error with FromChannelToDB",
				)
				continue
			}
			values = nil
		}
	}
}
