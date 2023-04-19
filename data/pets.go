package data

import (
	"encoding/json"
	"io"
	"time"
)

type Pet struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
	// Animal
	// Description string
	Sex       string `json:"gender"`
	Age       int    `json:"age"`
	KeptSince string `json:"Kept Since"`
}

func (p *Pet) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Pets []*Pet

func (p *Pets) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetPets() Pets {
	return PetList
}

var PetList = []*Pet{
	{
		ID:        1,
		Name:      "Mimi",
		Sex:       "Male",
		Age:       1,
		KeptSince: time.Date(2022, 5, 28, 0, 0, 0, 0, time.UTC).UTC().String(),
	},
}
