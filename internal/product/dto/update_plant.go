package dto

import (
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/google/uuid"
)

type UpdatePlantInput struct {
	Name string `json:"name"`
}

type UpdatePlantOutput struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
}

func (i *UpdatePlantInput) Validate() error {
	if i.Name == "" {
		return entity.ErrNameInvalid
	}

	return nil
}
