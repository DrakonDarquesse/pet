package handlers

import (
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
	r.Post("/{id:[0-9]+}", p.UpdatePet)
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
	p.l.Println("Handle POST Pet")

	// fetch the pets from datastore
	pet := &data.Pet{}

	//  decode data
	err := pet.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
	}
	p.l.Printf("Pet: %#v", pet)
	data.AddPet(pet)
}

func (p Pets) UpdatePet(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Pet Update")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	// fetch the pets from datastore
	pet := &data.Pet{}

	//  decode data
	err := pet.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
	}

	p.l.Printf("Pet: %#v", pet)
	err = data.UpdatePet(id, pet)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}
