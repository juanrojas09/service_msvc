package usecases

import (
	"context"
	"log"

	"github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
)

type (
	SaveServiceReviewsImpl struct {
		repo repositories.ServiceRepository
		log  *log.Logger
	}
)

func NewSaveServiceReviewsImpl(repo repositories.ServiceRepository, logger *log.Logger) interfaces.UseCases {
	return &SaveServiceReviewsImpl{
		repo: repo,
		log:  logger,
	}
}

func (s *SaveServiceReviewsImpl) Handle(ctx context.Context, params ...interface{}) (interface{}, error) {
	s.log.Println("Saving service review...")
	req := params[0].(repositories.SaveServiceReviewRequestDto)
	err := s.repo.SaveServiceReview(ctx, req)
	if err != nil {
		s.log.Printf("Error saving service review: %v", err)
		return nil, response.InternalServerError(err.Error())
	}
	return response.Created("Service review saved successfully", nil, nil), nil
}
