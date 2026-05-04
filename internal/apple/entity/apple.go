package entity

import (
	"github.com/google/uuid"
)

const (
	StatusNew   = "new"
	StatusError = "error"
)

type Apple struct {
	ID     uuid.UUID
	Name   string
	Status string
	Stuffs []string
}
