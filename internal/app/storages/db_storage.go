package storages

import (
	"context"
	"database/sql"
	"errors"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/structs"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

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

func (dbs *DBStorage) MakeDB() {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users(user_id text, login text, password text)`)
	if err != nil {
		panic(err)
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS out_points(order_id text, user_id text, summ numeric, processed_at timestamp with time zone)`)
	if err != nil {
		panic(err)
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS orders(order_id text, user_id text, status text, accrual numeric, uploaded_at timestamp with time zone)`)
	if err != nil {
		panic(err)
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS balances(user_id text, balance numeric, withdrawn numeric)`)
	if err != nil {
		panic(err)
	}
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

				_, err = db.ExecContext(ctx,
					`INSERT INTO users (user_id, login, password) values ($1, $2, $3)`,
					cookie, login, password)
				if err != nil {
					return false, "", err
				}

				_, err = db.ExecContext(ctx,
					`INSERT INTO balances (user_id, balance, withdrawn) values ($1, $2, $3)`,
					cookie, 0, 0)
				if err != nil {
					return false, "", err
				}

				return false, cookie, nil
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

func (dbs *DBStorage) CheckOrder(order string, userId string) (bool, bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, false, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uid string
	err = db.QueryRowContext(ctx,
		`SELECT user_id FROM orders WHERE order_id = $1`, order).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.ExecContext(ctx,
				`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
				order, userId, "NEW", 0, time.Now().Format(time.RFC3339))
			if err != nil {
				return false, false, err
			}

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

func (dbs *DBStorage) GetOrders(userId string) (bool, []structs.Order, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, status, accrual, uploaded_at FROM orders WHERE user_id = $1 order by uploaded_at desc`,
		userId)
	if err != nil {
		return false, nil, err
	}

	var items []structs.Order
	for rows.Next() {
		var item structs.Order
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

func (dbs *DBStorage) GetBalance(userId string) (structs.Balance, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return structs.Balance{}, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item structs.Balance
	err = db.QueryRowContext(ctx,
		`SELECT balance, withdrawn FROM balances WHERE user_id = $1`, userId).Scan(&item.Current, &item.WithDrawn)
	if err != nil {
		return structs.Balance{}, err
	}
	return item, nil
}

func (dbs *DBStorage) CheckBalance(userId string, orderId string, orderSum float32) (bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item structs.Balance
	err = db.QueryRowContext(ctx,
		`SELECT balance, withdrawn FROM balances WHERE user_id = $1`, userId).Scan(&item.Current, &item.WithDrawn)
	if err != nil {
		return false, err
	}

	if item.Current < orderSum {
		return false, nil
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO out_points (order_id, user_id, summ, processed_at) values ($1, $2, $3, $4)`,
		orderId, userId, orderSum, time.Now().Format(time.RFC3339))
	if err != nil {
		return false, err
	}

	item.Current -= orderSum
	item.WithDrawn += orderSum

	_, err = db.ExecContext(ctx,
		`UPDATE balances SET balance = $1, withdrawn = $2 WHERE user_id = $3`,
		item.Current, item.WithDrawn, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Out_points

func (dbs *DBStorage) GetOutPoints(userId string) (bool, []structs.OutPoints, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, summ, processed_at FROM out_points WHERE user_id = $1 order by processed_at desc`,
		userId)
	if err != nil {
		return false, nil, err
	}

	var items []structs.OutPoints
	for rows.Next() {
		var item structs.OutPoints
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
