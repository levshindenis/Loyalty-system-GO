package storages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/tools"
)

type DBStorage struct {
	address string
}

func (dbs *DBStorage) GetAddress() string {
	return dbs.address
}

func (dbs *DBStorage) SetAddress(value string) {
	dbs.address = value
}

//

func (dbs *DBStorage) MakeDB() error {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users(user_id text, login text, password text)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS orders(order_id text, user_id text, status text, accrual numeric, uploaded_at timestamp with time zone)`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS balances(user_id text, balance numeric, withdrawn numeric)`)
	if err != nil {
		return err
	}

	return nil
}

// User

func (dbs *DBStorage) CheckUser(login string, password string, param string) (bool, string, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, "", err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userId, pswrd string
	err = db.QueryRowContext(ctx,
		`SELECT user_id, password FROM users WHERE login = $1`, login).Scan(&userId, &pswrd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if param == "registration" {
				var counter int
				err = db.QueryRowContext(ctx, `SELECT count(*) FROM users`).Scan(&counter)
				if err != nil {
					return false, "", err
				}

				cookie, err := tools.GenerateCookie(counter + 1)
				if err != nil {
					return false, "", err
				}

				tx, err := db.Begin()
				if err != nil {
					return false, "", err
				}

				_, err = tx.ExecContext(ctx,
					`INSERT INTO users (user_id, login, password) values ($1, $2, $3)`,
					cookie, login, password)
				if err != nil {
					tx.Rollback()
					return false, "", err
				}

				_, err = tx.ExecContext(ctx,
					`INSERT INTO balances (user_id, balance, withdrawn) values ($1, $2, $3)`,
					cookie, 0, 0)
				if err != nil {
					tx.Rollback()
					return false, "", err
				}

				return false, cookie, tx.Commit()
			}

			return false, "", nil
		} else {
			return false, "", err
		}
	}

	if param == "login" {
		if pswrd != password {
			return false, "", nil
		}
	}

	return true, userId, nil
}

// Cookie

func (dbs *DBStorage) CheckCookie(cookie string) (bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userId string
	err = db.QueryRowContext(ctx,
		`SELECT user_id FROM users WHERE user_id = $1`, cookie).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// Order

func (dbs *DBStorage) CheckOrder(orderId string, userId string) (bool, bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, false, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uid string
	err = db.QueryRowContext(ctx,
		`SELECT user_id FROM orders WHERE order_id = $1`, orderId).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.ExecContext(ctx,
				`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
				orderId, userId, "NEW", nil, time.Now().Format(time.RFC3339))
			if err != nil {
				return false, false, err
			}
			uid = userId
			return false, false, nil
		} else {
			return false, false, err
		}
	}

	if uid == userId {
		return true, true, nil
	}
	return true, false, nil
}

func (dbs *DBStorage) GetOrders(userId string) (bool, []models.Order, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, status, accrual, uploaded_at FROM orders 
            WHERE user_id = $1 and status <> 'SOLD' order by uploaded_at desc`,
		userId)
	if err != nil {
		return false, nil, err
	}

	var items []models.Order
	for rows.Next() {
		var item models.Order
		if err = rows.Scan(&item.OrderId, &item.Status, &item.Accrual, &item.UploadedAt); err != nil {
			return false, nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return false, nil, err
	}

	if len(items) > 0 {
		return true, items, nil
	}

	return false, nil, nil
}

// Balance

func (dbs *DBStorage) GetBalance(userId string) (models.Balance, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return models.Balance{}, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item models.Balance
	err = db.QueryRowContext(ctx,
		`SELECT balance, withdrawn FROM balances WHERE user_id = $1`, userId).Scan(&item.Current, &item.WithDrawn)
	if err != nil {
		return models.Balance{}, err
	}
	return item, nil
}

func (dbs *DBStorage) CheckBalance(userId string, orderId string, orderSum float64) (bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var balance float64
	err = db.QueryRowContext(ctx,
		`SELECT balance FROM balances WHERE user_id = $1`, userId).Scan(&balance)
	if err != nil {
		return false, err
	}

	if balance < orderSum {
		return false, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
		orderId, userId, "SOLD", orderSum, time.Now().Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return false, err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE balances SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`,
		orderSum, userId)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit()
}

// Out_points

func (dbs *DBStorage) GetOutPoints(userId string) (bool, []models.OutPoints, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, accrual, uploaded_at FROM orders 
            WHERE user_id = $1 and status = 'SOLD' order by uploaded_at desc`,
		userId)
	if err != nil {
		return false, nil, err
	}

	var items []models.OutPoints
	for rows.Next() {
		var item models.OutPoints
		if err = rows.Scan(&item.OrderId, &item.Summ, &item.ProcessedAt); err != nil {
			return false, nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return false, nil, err
	}

	if len(items) > 0 {
		return true, items, nil
	}

	return false, nil, nil
}

// горутины

func (dbs *DBStorage) GetNewOrders() ([]models.Task, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, status, accrual FROM orders WHERE status in ('NEW', 'PROCESSING')`)
	if err != nil {
		return nil, err
	}

	var items []models.Task
	for rows.Next() {
		var item models.Task
		if err = rows.Scan(&item.OrderId, &item.Status, &item.Accrual); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (dbs *DBStorage) UpdateOrders(tasks []models.Task) error {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("tasks: ", tasks)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		_, err = tx.ExecContext(ctx,
			`UPDATE orders SET status = $1, accrual = $2, uploaded_at = $3 WHERE order_id = $4`,
			task.Status, *task.Accrual, time.Now().Format(time.RFC3339), task.OrderId)
		if err != nil {
			tx.Rollback()
			return err
		}

		var userId string
		err = tx.QueryRowContext(ctx,
			`SELECT user_id FROM orders WHERE order_id = $1`, task.OrderId).Scan(&userId)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE balances SET balance = balance + $1 WHERE user_id = $2`, *task.Accrual, userId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
