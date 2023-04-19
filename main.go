package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/drakondarquesse/pet/handlers"
)

func main() {

	// register handler function in the DefaultServeMux
	// http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {

	// })

	l := log.New(os.Stdout, "api", log.LstdFlags)
	hh := handlers.NewHello(l)

	sm := http.NewServeMux()

	sm.Handle("/hello", hh)

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
