package dto

import "github.com/google/uuid"

type CreatePlantOutput struct {
	ID uuid.UUID `json:"id"`
}

type CreatePlantInput struct {
	Name string `json:"name"`
}

func (i *CreatePlantInput) Validate() error {

	return nil
}
