package data

type Repository interface {
	All() ([]*Pet, error)
	AddPet(p *Pet) error
	UpdatePet(id int, p *Pet) error
	DeletePet(id int) error
}
