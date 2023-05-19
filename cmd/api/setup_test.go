package main

import (
	"log"
	"os"
	"testing"

	"github.com/drakondarquesse/pet/data"
)

var testPets *Pets

func TestMain(m *testing.M) {
	l := log.New(os.Stdout, "api", log.LstdFlags)
	testPets = NewPets(l, data.NewPetTestRepository(nil))

	os.Exit(m.Run())
}
