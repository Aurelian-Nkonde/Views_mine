package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"thousand.views_mine/cmd/handlers"
	"thousand.views_mine/internals/database/db_quaries"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	dns := os.Getenv("DNS")
	fmt.Println(dns)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dns)
	if err != nil {
		log.Fatal("error!")
	}
	defer conn.Close(ctx)

	var q *db_quaries.Queries = db_quaries.New(conn)

	h := handlers.App{
		Quaries: q,
		Ctx:     ctx,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/account", func(r chi.Router) {
		r.Get("/", h.GetAccounts)
		r.Get("/{id}", h.GetAccount)
		r.Put("/verify-email/{id}", h.VerifyEmail)
		r.Delete("/{id}", h.DeleteAccount)
		r.Post("/", h.CreateAccount)
	})
	r.Route("/view", func(r chi.Router) {
		r.Get("/", h.GetAllViews)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Not found!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Oops, 404!"))
	})

	s := &http.Server{
		Addr:           ":4500",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("App server is up and running...cheers!")
	if err := s.ListenAndServe(); err != nil {
		log.Panic("Error running the app!")
	}

}
