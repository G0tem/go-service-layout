package v1

import "github.com/G0tem/go-service-layout/internal/apple/usecase"

type Handlers struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handlers {
	return &Handlers{
		uc: uc,
	}
}
