package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type JsonUtil struct {
}

func (j *JsonUtil) ToJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(data)
}

func (j *JsonUtil) FromJSON(r io.Reader, data any) error {
	e := json.NewDecoder(r)
	return e.Decode(data)
}

type Animal struct {
}

type Pet struct {
	ID   int    `json:"-"`
	Name string `json:"name" validate:"required"`
	// Animal
	// Description string
	Sex       string `json:"gender" validate:"required"`
	Age       int    `json:"age,omitempty" validate:"required,gte=0,lte=100"`
	KeptSince string `json:"kept since,omitempty" validate:"required,date"`
}

// A custom PetModel that wraps sql.DB connection pool
type PetModel struct {
	DB *sql.DB
}

func (p *Pet) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Pet) ToJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(p)
}

func (p *Pet) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("date", validateDate)
	return validate.Struct(p)
}

func (m PetModel) All() (Pets, error) {
	rows, err := m.DB.Query("SELECT * FROM pets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pets []*Pet

	// Loop through rows, using Scan to assign column data to struct fields
	for rows.Next() {
		pet := &Pet{}
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Sex, &pet.Age, &pet.KeptSince); err != nil {
			return pets, err
		}
		pets = append(pets, pet)
	}
	if err = rows.Err(); err != nil {
		return pets, err
	}

	return pets, nil
}

func validateDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("01/02/2006", fl.Field().String())
	return err != nil
}

func (m PetModel) AddPet(p *Pet) error {
	var sqlStatement = `
	INSERT INTO pets (name, gender, age, keptsince)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	id := 0
	err := m.DB.QueryRow(sqlStatement, p.Name, p.Sex, p.Age, p.KeptSince).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (m PetModel) UpdatePet(id int, p *Pet) error {
	var sqlStatement = `
	UPDATE pets
	SET name = $1,
		gender = $2,
		age = $3,
		keptsince = $4 
	WHERE id = $5`

	result, err := m.DB.Exec(sqlStatement, p.Name, p.Sex, p.Age, p.KeptSince, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New("No data is updated")
	}
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

func (p *Pets) ToJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
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
		Age:       2,
		KeptSince: time.Date(2022, 5, 28, 0, 0, 0, 0, time.UTC).UTC().String(),
	},
}
