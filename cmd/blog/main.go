package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v7"

	"github.com/jadoint/micro/pkg/blog"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/db"
	"github.com/jadoint/micro/pkg/env"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/visitor"
)

func main() {
	// Logging
	file, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err.Error())
	}
	defer file.Close()
	log.SetOutput(file)

	// Load environment variables if
	// not already set.
	if os.Getenv("LISTEN") == "" {
		err = env.Load()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Database
	dbClient, err := db.GetClient()
	if err != nil {
		logger.Panic(err.Error())
	}
	defer dbClient.Master.Close()
	defer dbClient.Read.Close()

	// Cache
	redisClient := redis.NewClient(&redis.Options{Addr: os.Getenv("CACHE_ADDR")})
	defer redisClient.Close()

	// Clients
	clients := &conn.Clients{
		DB:    dbClient,
		Cache: redisClient,
	}

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
	// Rate limiter: first argument is "x requests / second" per IP
	lmt := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	r.Use(tollbooth_chi.LimitHandler(lmt))
	r.Use(visitor.Middleware)

	startPath := fmt.Sprintf(`/%s/`, os.Getenv("START_PATH"))
	r.Mount(startPath+"blog/tag", blog.RouteTag(clients))
	r.Mount(startPath+"blog", blog.RouteBlog(clients))

	srv := &http.Server{
		Addr:         os.Getenv("LISTEN"),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println(srv.ListenAndServeTLS(os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")))
}
