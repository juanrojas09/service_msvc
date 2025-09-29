package controllers

import (
	"context"

	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
)

type (
	Controller func(ctx context.Context, req interface{}) (interface{}, error)

	Endpoints struct {
		CreateServiceRequest          Controller
		GetServiceListByUserIdRequest Controller
		GetServiceDetailByIdRequest   Controller
		SaveServiceEvidence           Controller
		SaveServiceReviews            Controller
		GetProfessionalServiceList    Controller
	}

	UseCaseRegistry struct {
		CreateServiceRequestUseCase interfaces.UseCases
		ListServiceByUserIdUseCase  interfaces.UseCases
		GetServiceDetailByIdUseCase interfaces.UseCases
		SaveServiceEvidenceUseCase  interfaces.UseCases
		SaveServiceReviewsUseCase   interfaces.UseCases
		GetProfessionalServiceList  interfaces.UseCases
	}
)

func MakeEndpoints(ucr *UseCaseRegistry) Endpoints {
	return Endpoints{
		CreateServiceRequest:          makeCreateServiceRequestEndpoint(ucr.CreateServiceRequestUseCase),
		GetServiceListByUserIdRequest: makeGetServiceListByUserIdEndpoint(ucr.ListServiceByUserIdUseCase),
		GetServiceDetailByIdRequest:   makeGetServiceDetailByIdEndpoint(ucr.GetServiceDetailByIdUseCase),
		SaveServiceEvidence:           makeSaveServiceEvidenceEndpoint(ucr.SaveServiceEvidenceUseCase),
		SaveServiceReviews:            makeSaveServiceReviewsEndpoint(ucr.SaveServiceReviewsUseCase),
		GetProfessionalServiceList:    makeGetProfessionalServiceListEndpoint(ucr.GetProfessionalServiceList),
	}
}

func makeCreateServiceRequestEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}

func makeGetServiceListByUserIdEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}

func makeGetServiceDetailByIdEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}

func makeSaveServiceEvidenceEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}

func makeSaveServiceReviewsEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}

func makeGetProfessionalServiceListEndpoint(uc interfaces.UseCases) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return uc.Handle(ctx, req)
	}
}
