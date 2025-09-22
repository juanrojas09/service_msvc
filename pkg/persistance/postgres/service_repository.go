package postgres

import (
	"context"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/juanrojas09/core_domain/domain"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
	"gorm.io/gorm"
)

type (
	ServiceRepositoryImp struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewServiceRepository(db *gorm.DB, logger *log.Logger) repositories.ServiceRepository {
	return &ServiceRepositoryImp{
		db:  db,
		log: logger,
	}
}
func (s *ServiceRepositoryImp) CreateService(ctx context.Context, dto repositories.CreateServiceRequestDTO) (repositories.CreateServiceResponseDto, error) {
	s.log.Println("Inserting service into the database...")
	var resp repositories.CreateServiceResponseDto

	var statusRes domain.Status
	res := s.db.WithContext(ctx).Model(&domain.Status{}).Where("name = ?", "PENDIENTE").First(&statusRes)

	if res.Error == gorm.ErrRecordNotFound {
		return repositories.CreateServiceResponseDto{}, gorm.ErrRecordNotFound
	}

	var categoryRes domain.Categories
	res = s.db.WithContext(ctx).Model(&domain.Categories{}).Where("id = ?", dto.CategoryID).First(&categoryRes)

	if res.Error == gorm.ErrRecordNotFound {
		return repositories.CreateServiceResponseDto{}, gorm.ErrRecordNotFound
	}
	clientLng, err := strconv.ParseFloat(dto.ClientLongitude, 64)
	clientLat, err := strconv.ParseFloat(dto.ClientLatitude, 64)

	if err != nil {
		return repositories.CreateServiceResponseDto{}, err
	}

	tx := s.db.WithContext(ctx).Model(&domain.ServicesRequests{}).Create(
		&domain.ServicesRequests{
			ID:             uuid.New().String(),
			ProfessionalID: dto.ProfessionalID,
			ClientID:       dto.ClientID,
			Description:    dto.Description,
			Category:       categoryRes,
			LastClientLat:  &clientLat,
			LastClientLng:  &clientLng,
			Status:         statusRes,
		},
	).Select("id, client_id, professional_id, description, last_client_lat as client_latitude, last_client_lng as client_longitude").Scan(&resp)

	if tx.Error != nil {
		s.log.Printf("Error inserting service into the database: %v", tx.Error)
		return repositories.CreateServiceResponseDto{}, tx.Error
	}

	resp.CategoryName = categoryRes.Name
	resp.Status = string(statusRes.Name)
	return resp, nil

}
