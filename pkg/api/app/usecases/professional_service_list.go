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
	ProfessionalListServiceImpl struct {
		repo repositories.ServiceRepository
		log  *log.Logger
	}
)

func NewProfessionalListServiceImpl(repo repositories.ServiceRepository, logger *log.Logger) interfaces.UseCases {
	return &ProfessionalListServiceImpl{
		repo: repo,
		log:  logger,
	}
}
func (p *ProfessionalListServiceImpl) Handle(ctx context.Context, params ...interface{}) (interface{}, error) {
	r := params[0].(repositories.ProfessionalServiceListRequestDto)

	count, err := p.repo.GetProfessionalServicesCount(ctx, r.ProfessionalID)
	if err != nil {
		p.log.Println("Error counting professional services:", err)
		return nil, response.InternalServerError(err.Error())
	}
	meta, err := meta.New(count, r.Page, r.Limit, "10")
	if err != nil {
		p.log.Println("Error creating metadata:", err)
		return nil, response.BadRequest(err.Error())
	}
	res, err := p.repo.GetProfessionalServicesById(ctx, r.ProfessionalID, meta.Offset(), meta.Limit())
	if err != nil {
		p.log.Println("Error getting professional services:", err)
		return nil, response.InternalServerError(err.Error())
	}
	return response.Ok("Professional Services Fetched Successfully", res, meta), nil
}
