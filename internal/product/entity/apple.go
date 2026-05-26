package entity

import (
	"github.com/google/uuid"
)

const (
	StatusNew    = "new"
	StatusUpdate = "update"
	StatusError  = "error"
)

type Apple struct {
	ID      uuid.UUID
	PlantID uuid.UUID
	Name    string
	Status  string
	Stuffs  []string
}
