package main

import (
	"log"
	"realTime/server/db"
)

func main() {
	_, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

}
