package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Hello satisfies the Handler interface because it has
// ServeHTTP method
type Hello struct {
	l *log.Logger
}

// constructor
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// method of type Hello
func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	// read the data from body
	d, err := io.ReadAll(r.Body)
	h.l.Printf("Hello, %s", d)

	// handle error
	if err != nil {
		http.Error(w, "Oops", http.StatusBadRequest)
		return
	}

	// write to ResponseWriter
	fmt.Fprintf(w, "Hello, %s\n", d)
	w.Write([]byte("Thank You"))
}
