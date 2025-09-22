package controllers

import (
	"context"

	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
)

type (
	Controller func(ctx context.Context, req interface{}) (interface{}, error)

	Endpoints struct {
		CreateServiceRequest Controller
	}

	UseCaseRegistry struct {
		CreateServiceRequestUseCase interfaces.UseCases
	}
)

func MakeEndpoints(ucr *UseCaseRegistry) Endpoints {
	return Endpoints{
		CreateServiceRequest: makeCreateServiceRequestEndpoint(ucr.CreateServiceRequestUseCase),
	}
}

func makeCreateServiceRequestEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}
