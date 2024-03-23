package database

import (
	"context"
	"time"
)

func (dbs *DBStorage) MakeDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := dbs.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users(user_id text, login text, password text)`)
	if err != nil {
		return err
	}

	_, err = dbs.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS orders(order_id text, user_id text, status text, accrual numeric, uploaded_at timestamp with time zone)`)
	if err != nil {
		return err
	}

	_, err = dbs.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS balances(user_id text, balance numeric, withdrawn numeric)`)
	if err != nil {
		return err
	}

	return nil
}
