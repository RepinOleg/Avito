package dbs

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewDB(cfg DBConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)

	connect, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	err = connect.Ping()
	if err != nil {
		return nil, err
	}

	return connect, nil
}
