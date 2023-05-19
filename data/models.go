package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Animal struct {
	Sex string `json:"gender" validate:"required"`
	Age int    `json:"age,omitempty" validate:"required,gte=0,lte=100"`
}

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	Animal
	KeptSince string `json:"kept since,omitempty" validate:"required,date"`
}

// A custom PetModel that wraps sql.DB connection pool
type PetRepository struct {
	DB *sql.DB
}

func NewPetRepository(db *sql.DB) PetRepository {
	return PetRepository{db}
}

func (m PetRepository) All() ([]*Pet, error) {
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

func (m PetRepository) AddPet(p *Pet) error {
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

func (m PetRepository) UpdatePet(id int, p *Pet) error {
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
		return errors.New("No pet is updated")
	}
	return nil
}

func (m PetRepository) DeletePet(id int) error {
	var sqlStatement = `
	DELETE FROM pets
	WHERE id = $1`

	result, err := m.DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New("No pet is deleted")
	}
	return nil
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
