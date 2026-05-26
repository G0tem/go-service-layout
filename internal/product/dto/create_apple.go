package dto

import (
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/google/uuid"
)

type CreateAppleOutput struct {
	ID uuid.UUID `json:"id"`
}

type CreateAppleInput struct {
	Name    string    `json:"name"`
	PlantID uuid.UUID `json:"plant_id"`
}

func (i *CreateAppleInput) Validate() error {
	if i.Name == "" {
		return entity.ErrNameInvalid
	}

	if i.PlantID == uuid.Nil {
		return entity.ErrUUIDInvalid
	}

	return nil
}
