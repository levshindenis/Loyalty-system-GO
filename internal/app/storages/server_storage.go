package storages

import (
	"github.com/levshindenis/Loyalty-system-GO/internal/app/config"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

type ServerStorage struct {
	sc config.ServerConfig
	ds DBStorage
}

func (serv *ServerStorage) ParseFlags() {
	serv.sc.ParseFlags()
}

func (serv *ServerStorage) Init() {
	serv.ds = DBStorage{address: serv.sc.GetDbURI()}
	serv.ds.MakeDB()
}

func (serv *ServerStorage) GetRunAddress() string {
	return serv.sc.GetRunAddress()
}

func (serv *ServerStorage) CheckUser(login string, password string, param string) (bool, string, error) {
	return serv.ds.CheckUser(login, password, param)
}

func (serv *ServerStorage) CheckCookie(cookie string) (bool, error) {
	return serv.ds.CheckCookie(cookie)
}

func (serv *ServerStorage) CheckOrder(orderId string, userId string, param string) (bool, bool, error) {
	return serv.ds.CheckOrder(orderId, userId, param)
}

func (serv *ServerStorage) GetOrders(userId string) (bool, []models.Order, error) {
	return serv.ds.GetOrders(userId)
}

func (serv *ServerStorage) GetBalance(userId string) (models.Balance, error) {
	return serv.ds.GetBalance(userId)
}

func (serv *ServerStorage) CheckBalance(userId string, orderId string, orderSum float32) (bool, error) {
	return serv.ds.CheckBalance(userId, orderId, orderSum)
}

func (serv *ServerStorage) GetOutPoints(userId string) (bool, []models.OutPoints, error) {
	return serv.ds.GetOutPoints(userId)
}

func (serv *ServerStorage) SetDBAddress(value string) {
	serv.ds.SetAddress(value)
}
