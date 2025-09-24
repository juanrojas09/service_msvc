package usecases

import (
	"context"
	"log"

	"github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/gocourse_meta/meta"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
)

type (
	ServiceListByUserIdImpl struct {
		repo repositories.ServiceRepository
		log  *log.Logger
	}
)

func NewServiceListByUserIdImpl(repo repositories.ServiceRepository, logger *log.Logger) interfaces.UseCases {
	return &ServiceListByUserIdImpl{
		repo: repo,
		log:  logger,
	}
}

func (s *ServiceListByUserIdImpl) Handle(ctx context.Context, params ...interface{}) (interface{}, error) {
	s.log.Println("Listing services by user ID...")
	req := params[0].(repositories.ServiceListRequestDTO)
	count, err := s.repo.CountServicesByUserId(ctx, req.UserID)
	if err != nil {
		s.log.Println("Error counting services:", err)
		return nil, response.InternalServerError(err.Error())
	}
	meta, err := meta.New(count, req.Page, req.Limit, "10")
	if err != nil {
		s.log.Println("Error creating metadata:", err)
		return nil, response.BadRequest(err.Error())
	}

	res, err := s.repo.GetServicesByUserId(ctx, req.UserID, meta.Offset(), meta.Limit())
	if err != nil {
		s.log.Println("Error getting services:", err)
		return nil, response.InternalServerError(err.Error())
	}

	return response.Ok("Services Fetched Successfully", repositories.ServiceListResponseDTO{
		Data: res,
	}, meta), nil
}
