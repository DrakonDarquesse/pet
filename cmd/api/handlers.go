package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/drakondarquesse/pet/data"
	"github.com/go-chi/chi/v5"
)

type Pets struct {
	l          *log.Logger
	repository data.Repository
	jsonUtil   interface {
		ToJSON(w http.ResponseWriter, data any) error
		FromJSON(r io.Reader, data any) error
	}
}

// Initialize Pets with logger and Repository
func NewPets(l *log.Logger, repository data.Repository) *Pets {
	return &Pets{l, repository, &data.JsonUtil{}}
}

func (p Pets) MountRoutes(r chi.Router) {
	r.Get("/", p.GetPets)
	r.Post("/", p.AddPet)
	r.Put("/{id:[0-9]+}", p.UpdatePet)
	r.Delete("/{id:[0-9]+}", p.DeletePet)
}

func (p Pets) GetPets(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Pets")

	pets, err := p.repository.All()

	if err != nil {
		p.l.Printf("Error: %#v", err)
	}

	// serialize to JSON
	err = p.jsonUtil.ToJSON(w, pets)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p Pets) AddPet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Pet Add")

	pet := &data.Pet{}

	err := p.jsonUtil.FromJSON(r.Body, pet)

	if err != nil {
		http.Error(w, "Invalid Input Format", http.StatusBadRequest)
		return
	}

	// validate data
	err = pet.Validate()
	if err != nil {
		p.l.Printf("Error validating pet %#v", err)
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
		return
	}

	p.l.Printf("Pet: %#v", pet)
	p.repository.AddPet(pet)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Create Success"))
}

func (p Pets) UpdatePet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Pet Update")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	pet := &data.Pet{}

	//  decode data
	err := p.jsonUtil.FromJSON(r.Body, pet)
	if err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
	}

	// validate data
	err = pet.Validate()
	if err != nil {
		p.l.Printf("Error validating pet %#v", err)
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
		return
	}

	err = p.repository.UpdatePet(id, pet)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Update Success"))
}

func (p Pets) DeletePet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Pet Delete")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := p.repository.DeletePet(id)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
	}
	w.Write([]byte("Delete Success"))
}
