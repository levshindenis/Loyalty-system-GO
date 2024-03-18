package server

import "github.com/levshindenis/Loyalty-system-GO/internal/app/models"

func (serv *Storage) CheckUser(login string, password string, param string) (bool, string, error) {
	return serv.dbs.GetData().CheckUser(login, password, param)
}

func (serv *Storage) CheckUserCookie(cookie string) (bool, error) {
	return serv.dbs.GetData().CheckUserCookie(cookie)
}

func (serv *Storage) CheckUserOrder(orderID string, userID string) (bool, bool, error) {
	return serv.dbs.GetData().CheckUserOrder(orderID, userID)
}

func (serv *Storage) GetUserOrders(userID string) (bool, []models.Order, error) {
	return serv.dbs.GetData().GetUserOrders(userID)
}

func (serv *Storage) GetUserBalance(userID string) (models.Balance, error) {
	return serv.dbs.GetData().GetUserBalance(userID)
}

func (serv *Storage) CheckUserBalance(userID string, orderID string, orderSum float64) (bool, error) {
	return serv.dbs.GetData().CheckUserBalance(userID, orderID, orderSum)
}

func (serv *Storage) GetUserOutPoints(userID string) (bool, []models.OutPoints, error) {
	return serv.dbs.GetData().GetUserOutPoints(userID)
}
