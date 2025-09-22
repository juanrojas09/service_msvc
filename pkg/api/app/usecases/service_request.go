package usecases

import (
	"context"
	"log"

	r "github.com/juanrojas09/go_lib_response/response"

	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
)

type (
	ServiceRequestImpl struct {
		repo repositories.ServiceRepository
		log  *log.Logger
	}
)

func NewServiceRequestImpl(repo repositories.ServiceRepository, logger *log.Logger) interfaces.UseCases {
	return &ServiceRequestImpl{
		repo: repo,
		log:  logger,
	}
}

func (s *ServiceRequestImpl) Handle(ctx context.Context, params ...interface{}) (interface{}, error) {
	s.log.Println("Creating service...")
	req := params[0].(repositories.CreateServiceRequestDTO)
	isValid, err := s.repo.ValidateExistingPendingServiceFromClientToProfessional(ctx, req.ClientID, req.ProfessionalID)
	if err != nil {
		s.log.Printf("Error validating existing pending service: %v", err)
		return repositories.CreateServiceResponseDto{}, r.InternalServerError(err.Error())
	}
	if isValid {
		s.log.Println("Existing pending service found.")
		return repositories.CreateServiceResponseDto{}, r.BadRequest("Existing pending service found.")
	}
	response, err := s.repo.CreateService(ctx, req)
	if err != nil {
		s.log.Printf("Error creating service: %v", err)
		return repositories.CreateServiceResponseDto{}, r.InternalServerError(err.Error())
	}
	return r.Created("Service Created Successfully", response, nil), nil
}
