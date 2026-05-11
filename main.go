package main

import (
	"log"
	"todo-calendar/db"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateTables(conn)
	if err != nil {
		log.Fatal(err)
	}
}
