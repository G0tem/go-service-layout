package dto

import "github.com/google/uuid"

type CreatePineAppleOutput struct {
	ID uuid.UUID `json:"id"`
}

type CreatePineAppleInput struct {
	Name string `json:"name"`
}

func (i *CreatePineAppleInput) Validate() error {

	return nil
}
