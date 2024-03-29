package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/db"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/user"
	"github.com/jadoint/micro/pkg/visitor"
)

func main() {
	// Load environment variables if
	// not already set.
	if os.Getenv("LISTEN") == "" {
		log.Fatal("LISTEN is not set")
	}

	// Database
	dbClient, err := db.GetClient()
	if err != nil {
		logger.Panic(err.Error())
	}
	defer dbClient.Master.Close()
	defer dbClient.Read.Close()

	// Clients
	clients := &conn.Clients{DB: dbClient}

	// Routes
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("SITE_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	if os.Getenv("ENV") == "development" {
		r.Use(middleware.SetHeader("Access-Control-Allow-Origin", os.Getenv("SITE_URL")))
	}
	r.Use(visitor.Middleware)

	startPath := fmt.Sprintf(`/%s/`, os.Getenv("START_PATH"))
	r.Mount(startPath+"auth", user.RouteAuth(clients))
	r.Mount(startPath+"user", user.RouteUser(clients))

	srv := &http.Server{
		Addr:         os.Getenv("LISTEN"),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println(srv.ListenAndServeTLS(os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")))
}
