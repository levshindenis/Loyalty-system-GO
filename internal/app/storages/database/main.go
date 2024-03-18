package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

type DBFuncs interface {
	CheckUser(login string, password string, param string) (bool, string, error)
	CheckUserCookie(cookie string) (bool, error)
	CheckUserOrder(orderID string, userID string) (bool, bool, error)
	GetUserOrders(userID string) (bool, []models.Order, error)
	GetUserBalance(db *sql.DB, userID string) (models.Balance, error)
	CheckUserBalance(userID string, orderID string, orderSum float64) (bool, error)
	GetUserOutPoints(userID string) (bool, []models.OutPoints, error)
	GetNewOrders() ([]models.Task, error)
	UpdateOrders(tasks []models.Task) error
}

type DBStorage struct {
	Data DBFuncs
}
