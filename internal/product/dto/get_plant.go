package dto

import "github.com/google/uuid"

type GetPlantOutput struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
}

type GetPlantInput struct {
	ID uuid.UUID
}

func (i *GetPlantInput) Validate() error {

	return nil
}
