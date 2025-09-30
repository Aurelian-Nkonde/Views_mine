package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"thousand.views_mine/cmd/handlers"
	"thousand.views_mine/internals/database/db_quaries"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	tokenAuth = jwtauth.New("HS256", []byte("secret-code"), nil)

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
		Token:   tokenAuth,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/account", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Get("/", h.GetAccounts)
			r.Get("/{id}", h.GetAccount)
			r.Delete("/{id}", h.DeleteAccount)
		})
		r.Put("/verify-email/{id}", h.VerifyEmail)
		r.Post("/signup", h.CreateAccount)
		r.Post("/login", h.Login)
	})
	r.Route("/view", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Get("/user/{id}", h.GetUserViews)
			r.Post("/", h.CreateView)
			r.Delete("/{id}", h.DeleteView)
		})
		r.Get("/", h.GetAllViews)
		r.Get("/{id}", h.GetView)
		r.Get("/public", h.GetAllPublicViews)
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
