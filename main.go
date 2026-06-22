package main

import (
	"log"
	"net/http"

	"todo-calendar/auth"
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

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/landing.html")
	})
	http.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/app.html")
	})
	http.HandleFunc("GET /stats", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/stats.html")
	})
	http.HandleFunc("GET /archive", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/archive.html")
	})
	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	http.HandleFunc("GET /friends", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/friends.html")
	})
	http.HandleFunc("GET /shared", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/shared.html")
	})

	http.HandleFunc("POST /api/register", h.Register)
	http.HandleFunc("POST /api/login", h.Login)
	http.HandleFunc("POST /api/logout", h.Logout)
	http.HandleFunc("GET /api/me", h.Me)

	http.HandleFunc("GET /tasks", auth.RequireAuth(conn, h.GetTasks))
	http.HandleFunc("POST /tasks", auth.RequireAuth(conn, h.CreateTask))
	http.HandleFunc("PUT /tasks/{id}", auth.RequireAuth(conn, h.UpdateTask))
	http.HandleFunc("DELETE /tasks/{id}", auth.RequireAuth(conn, h.DeleteTask))

	http.HandleFunc("POST /api/friends/requests", auth.RequireAuth(conn, h.SendFriendRequest))
	http.HandleFunc("GET /api/friends/requests", auth.RequireAuth(conn, h.ListIncomingRequests))
	http.HandleFunc("POST /api/friends/accept", auth.RequireAuth(conn, h.AcceptFriendRequest))
	http.HandleFunc("POST /api/friends/decline", auth.RequireAuth(conn, h.DeclineFriendRequest))
	http.HandleFunc("GET /api/friends", auth.RequireAuth(conn, h.ListFriends))

	http.HandleFunc("POST /api/access/grant", auth.RequireAuth(conn, h.GrantAccess))
	http.HandleFunc("POST /api/access/revoke", auth.RequireAuth(conn, h.RevokeAccess))
	http.HandleFunc("GET /api/access/shared-with-me", auth.RequireAuth(conn, h.ListSharedWithMe))
	http.HandleFunc("GET /api/access/my-grants", auth.RequireAuth(conn, h.ListMyGrants))
	http.HandleFunc("GET /api/access/calendar/{ownerID}", auth.RequireAuth(conn, h.GetFriendTasks))

	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
