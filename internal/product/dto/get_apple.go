package dto

import "github.com/google/uuid"

type GetAppleOutput struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	PlantID uuid.UUID `json:"plant_id"`
}

type GetAppleInput struct {
	ID uuid.UUID
}

func (i *GetAppleInput) Validate() error {

	return nil
}
