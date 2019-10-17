package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/db"
	"github.com/jadoint/micro/pkg/env"
	"github.com/jadoint/micro/pkg/logger"
	appmiddleware "github.com/jadoint/micro/pkg/middleware"
	"github.com/jadoint/micro/pkg/user/route"
)

func main() {
	// Logging
	file, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	// Environment variables
	err = env.Load()
	if err != nil {
		log.Fatal(err.Error())
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
	r.Use(appmiddleware.Middleware)

	startPath := fmt.Sprintf(`/%s/`, os.Getenv("START_PATH"))
	r.Mount(startPath+"auth", route.AuthRouter(clients))
	r.Mount(startPath+"user", route.UserRouter(clients))

	srv := &http.Server{
		Addr:         os.Getenv("LISTEN"),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println(srv.ListenAndServeTLS(os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")))
}
