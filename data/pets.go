package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

type Pet struct {
	ID   int    `json:"-"`
	Name string `json:"name" validate:"required"`
	// Animal
	// Description string
	Sex       string `json:"gender" validate:"required"`
	Age       int    `json:"age" validate:"gte=0, lte=100"`
	KeptSince string `json:"kept since" validate:"date"`
}

func (p *Pet) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Pet) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("date", validateDate)
	return validate.Struct(p)
}

func validateDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("01/02/2006", fl.Field().String())
	return err != nil
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

func DeletePet(id int) error {
	_, err := FindPet(id)
	if err != nil {
		return err
	}

	// perform delete

	return nil
}

func FindPet(id int) (int, error) {
	for i, p := range petList {
		if p.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("pet not found")
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

// sqlStatement := `
// INSERT INTO pets (name, gender, age, keptsince)
// VALUES ($1, $2, $3, $4)
// RETURNING id`
// id := 0
// err = db.QueryRow(sqlStatement, "Mimi", "m", 1, time.Date(2022, 5, 28, 0, 0, 0, 0, time.UTC).UTC().String()).Scan(&id)
// if err != nil {
// 	panic(err)
// }
// fmt.Println("New record ID is:", id)
