package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/drakondarquesse/pet/data"
	"github.com/go-chi/chi"
)

type Pets struct {
	l *log.Logger
}

func NewPets(l *log.Logger) *Pets {
	return &Pets{l}
}

func (p Pets) MountRoutes(r chi.Router) {
	r.Get("/", p.GetPets)
	r.Post("/add", p.AddPet)
	r.Put("/{id:[0-9]+}", p.UpdatePet)
}

func (p Pets) GetPets(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Pets")

	// fetch the pets from datastore
	petList := data.GetPets()

	// serialize to JSON
	err := petList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p Pets) AddPet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Pet Add")

	pet := &data.Pet{}

	//  decode data
	err := pet.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
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
	data.AddPet(pet)
	w.Write([]byte("Create Success"))

}

func (p Pets) UpdatePet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Pet Update")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	pet := &data.Pet{}

	//  decode data
	err := pet.FromJSON(r.Body)
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

	p.l.Printf("Pet: %#v", pet)
	err = data.UpdatePet(id, pet)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
	}
}

func (p Pets) DeletePet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Pet Delete")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := data.DeletePet(id)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
	}
}
