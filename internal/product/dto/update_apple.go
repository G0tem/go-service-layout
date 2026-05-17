package dto

import (
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/google/uuid"
)

type UpdateAppleInput struct {
	Name string `json:"name"`
}

type UpdateAppleOutput struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

func (i *UpdateAppleInput) Validate() error {
	if i.Name == "" {
		return entity.ErrNameInvalid
	}

	return nil
}
