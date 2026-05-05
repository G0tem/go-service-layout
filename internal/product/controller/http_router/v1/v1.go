package v1

import "github.com/G0tem/go-service-layout/internal/product/usecase"

type Handlers struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handlers {
	return &Handlers{
		uc: uc,
	}
}
