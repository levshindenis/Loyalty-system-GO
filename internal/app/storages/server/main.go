package server

import (
	"context"
	"database/sql"
	"go.uber.org/zap"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/accrual"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/storages/database"
)

type Storage struct {
	dbs    database.DBStorage
	db     *sql.DB
	sl     zap.SugaredLogger
	fromDB chan models.Task
	toDB   chan models.Task
	ctx    context.Context
	cancel context.CancelFunc
}

func (serv *Storage) Init(db *sql.DB, AccSysAddr string) error {
	serv.db = db

	if err := serv.dbs.MakeDB(serv.db); err != nil {
		return err
	}

	serv.ctx, serv.cancel = context.WithCancel(context.Background())

	serv.NewLogger()

	serv.fromDB = make(chan models.Task, 1024)
	serv.toDB = make(chan models.Task, 1024)

	go serv.FromDBToChannel(serv.fromDB, serv.ctx)

	for i := 0; i < 5; i++ {
		w := accrual.NewCompareWorker(i, serv.fromDB, serv.toDB, AccSysAddr)
		go w.Loop(&serv.sl)
	}

	go serv.FromChannelToDB(serv.toDB, serv.ctx)

	return nil
}

func (serv *Storage) Terminate() {
	serv.cancel()
}

func (serv *Storage) NewLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	serv.sl = *logger.Sugar()
}
