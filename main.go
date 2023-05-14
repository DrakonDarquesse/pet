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

	"github.com/drakondarquesse/pet/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
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

	l := log.New(os.Stdout, "api", log.LstdFlags)
	hh := handlers.NewHello(l)
	ph := handlers.NewPets(l)

	sm.Handle("/hello", hh)

	sm.Route("/pets", func(r chi.Router) {
		ph.MountRoutes(r)
	})

	// create a server
	s := &http.Server{
		Addr:        ":9898",
		Handler:     sm,
		IdleTimeout: 120 * time.Second,
	}

	// run a server
	// http.ListenAndServe(":9898", sm)

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
