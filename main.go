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

	// Страницы
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/landing.html")
	})
	http.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/app.html")
	})

	// API
	http.HandleFunc("GET /tasks", h.GetTasks)
	http.HandleFunc("POST /tasks", h.CreateTask)
	http.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

	// CSS, JS и т.д.
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
