package db

import (
	"database/sql"
	"realTime/config"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

var Env = config.Env

func NewDatabase() (*Database, error) {
	Env = config.InitConfig("main")

	return databaseMysql(mysql.Config{
		User:   Env.DBUser,
		Passwd: Env.DBPassword,
		Net:    "tcp",

		Addr:                 Env.DBAddress,
		DBName:               Env.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	})

}

func databaseMysql(cfg mysql.Config) (*Database, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, err
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}
