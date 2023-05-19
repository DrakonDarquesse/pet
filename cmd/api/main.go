package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/drakondarquesse/pet/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

const (
	host     = "host.docker.internal"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "pet"
)

func main() {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// create a chi mux
	sm := chi.NewRouter()
	sm.Use(middleware.RequestID)
	sm.Use(middleware.Logger)
	sm.Use(middleware.Recoverer)

	// specify who is allowed to connect
	sm.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	sm.Use(middleware.Heartbeat("/ping"))

	l := log.New(os.Stdout, "api", log.LstdFlags)
	ph := NewPets(l, data.NewPetRepository(db))

	sm.Route("/pets", func(r chi.Router) {
		ph.MountRoutes(r)
	})

	// create a server
	s := &http.Server{
		Addr:        ":9898",
		Handler:     sm,
		IdleTimeout: 120 * time.Second,
	}

	// a goroutine to run the server so it doesn't block
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// create a channel to receive os.Signal
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
