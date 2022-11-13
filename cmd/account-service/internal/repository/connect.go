package repository

import (
	"database/sql"
	"fmt"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
)

func Connect(dbConfig env.DBConfig) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
