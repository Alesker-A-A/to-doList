package main

import (
	"log"
	"net/http"
	"todo-calendar/db"
	"todo-calendar/handlers"
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

	h := handlers.Handler{DB: conn}
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("GET /tasks", h.GetTasks)
	http.HandleFunc("POST /tasks", h.CreateTask)
	http.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
