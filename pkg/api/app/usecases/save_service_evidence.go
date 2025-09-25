package usecases

import (
	"context"
	"log"

	"github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
)

type (
	SaveServiceEvidenceImpl struct {
		repo repositories.ServiceRepository
		log  *log.Logger
	}
)

func NewSaveServiceEvidenceImpl(repo repositories.ServiceRepository, logger *log.Logger) interfaces.UseCases {
	return &SaveServiceEvidenceImpl{
		repo: repo,
		log:  logger,
	}
}

func (s *SaveServiceEvidenceImpl) Handle(ctx context.Context, params ...interface{}) (interface{}, error) {
	s.log.Println("Saving service evidence...")
	req := params[0].(repositories.SaveServiceEvidenceRequestDto)
	err := s.repo.SaveServiceEvidence(ctx, req)
	if err != nil {
		s.log.Printf("Error saving service evidence: %v", err)
		return nil, response.InternalServerError(err.Error())
	}
	return response.Created("Service evidence saved successfully", nil, nil), nil
}
