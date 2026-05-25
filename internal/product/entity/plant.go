package entity

import (
	"github.com/google/uuid"
)

type Name string

type Plant struct {
	ID     uuid.UUID
	Name   Name
	Status Status
	Stuffs []string
}

// New validate and create new Plant
func New(name, status string) (Plant, error) {
	if name == "" {
		return Plant{}, ErrNameInvalid
	}

	if status == "" {
		return Plant{}, ErrStatusInvalid
	}

	return Plant{
		ID:     uuid.New(),
		Name:   Name(name),
		Status: NewStatus(status),
	}, nil
}

func (a *Plant) AddStuff(stuff string) {
	a.Stuffs = append(a.Stuffs, stuff)
}

func (a *Plant) ChangeStatus(status Status) {
	a.Status = status
}

func (a *Plant) GetID() uuid.UUID {
	return a.ID
}

func (a *Plant) GetName() Name {
	return a.Name
}

func (a *Plant) GetStatus() Status {
	return a.Status
}

func (a *Plant) GetStuffs() []string {
	return a.Stuffs
}
