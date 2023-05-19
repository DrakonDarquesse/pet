package data

import (
	"database/sql"
)

type PetTestRepository struct {
	DB *sql.DB
}

func NewPetTestRepository(db *sql.DB) PetTestRepository {
	return PetTestRepository{nil}
}

func (m PetTestRepository) All() ([]*Pet, error) {
	var pets []*Pet

	return pets, nil
}

func (m PetTestRepository) AddPet(p *Pet) error {

	return nil
}

func (m PetTestRepository) UpdatePet(id int, p *Pet) error {

	return nil
}

func (m PetTestRepository) DeletePet(id int) error {

	return nil
}
