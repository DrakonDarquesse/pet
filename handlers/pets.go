package handlers

import "log"

type Pets struct {
	l *log.Logger
}

func NewPets(l *log.Logger) *Pets {
	return &Pets{l}
}
