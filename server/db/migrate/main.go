package main

import (
	"log"
	"os"
	"realTime/server/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	driver, err := mysql.WithInstance(db.DB(), &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://server/db/migrate", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		err := m.Up()
		if err != nil {
			log.Fatal("gagal membuat :", err)
		}
	}
	if cmd == "down" {
		err := m.Down()
		if err != nil {
			log.Fatal(err)
		}
	}
}
