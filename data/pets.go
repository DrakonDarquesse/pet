package data

import (
	"encoding/json"
	"fmt"
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

func AddPet(p *Pet) {
	petList = append(petList, p)
}

func UpdatePet(id int, p *Pet) error {
	pos, err := FindPet(id)
	if err != nil {
		return err
	}
	p.ID = id
	petList[pos] = p
	return nil
}

func FindPet(id int) (int, error) {
	for i, p := range petList {
		if p.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("product not found")
}

type Pets []*Pet

func (p *Pets) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetPets() Pets {
	return petList
}

var petList = []*Pet{
	{
		ID:        1,
		Name:      "Mimi",
		Sex:       "Male",
		Age:       1,
		KeptSince: time.Date(2022, 5, 28, 0, 0, 0, 0, time.UTC).UTC().String(),
	},
}
