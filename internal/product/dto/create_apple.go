package dto

import (
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/google/uuid"
)

type CreateAppleOutput struct {
	ID uuid.UUID `json:"id"`
}

type CreateAppleInput struct {
	Name string `json:"name"`
}

func (i *CreateAppleInput) Validate() error {
	if i.Name == "" {
		return entity.ErrNameInvalid
	}

	return nil
}
