package entity

import (
	"github.com/google/uuid"
)

type Name string

type PineApple struct {
	id     uuid.UUID
	name   Name
	status Status
	stuffs []string
}

// New validate and create new PineApple
func New(name, status string) (PineApple, error) {
	if name == "" {
		return PineApple{}, ErrNameInvalid
	}

	if status == "" {
		return PineApple{}, ErrStatusInvalid
	}

	return PineApple{
		id:     uuid.New(),
		name:   Name(name),
		status: NewStatus(status),
	}, nil
}

func (a *PineApple) AddStuff(stuff string) {
	a.stuffs = append(a.stuffs, stuff)
}

func (a *PineApple) ChangeStatus(status Status) {
	a.status = status
}

func (a *PineApple) GetID() uuid.UUID {
	return a.id
}

func (a *PineApple) GetName() Name {
	return a.name
}

func (a *PineApple) GetStatus() Status {
	return a.status
}

func (a *PineApple) GetStuffs() []string {
	return a.stuffs
}
