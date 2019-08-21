package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/db"
	"github.com/jadoint/micro/env"
	"github.com/jadoint/micro/logger"
	"github.com/jadoint/micro/user/route"
	"github.com/jadoint/micro/visitor"
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
	r.Use(visitor.Middleware)

	startPath := fmt.Sprintf(`/%s/`, os.Getenv("START_PATH"))
	r.Mount(startPath+"auth", route.AuthRouter(clients))

	srv := &http.Server{
		Addr:         os.Getenv("LISTEN"),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println(srv.ListenAndServe())
}
