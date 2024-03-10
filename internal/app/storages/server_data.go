package storages

import "github.com/levshindenis/Loyalty-system-GO/internal/app/models"

type BaseFuncs interface {
	CheckUser(login string, password string, param string) (bool, string, error)
	CheckCookie(cookie string) (bool, error)
	CheckOrder(orderID string, userID string) (bool, bool, error)
	GetOrders(userID string) (bool, []models.Order, error)
	GetBalance(userID string) (models.Balance, error)
	CheckBalance(userID string, orderID string, orderSum float64) (bool, error)
	GetOutPoints(userID string) (bool, []models.OutPoints, error)
	GetNewOrders() ([]models.Task, error)
	UpdateOrders(tasks []models.Task) error
}

type ServerData struct {
	data BaseFuncs
}
