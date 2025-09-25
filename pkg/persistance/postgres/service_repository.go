package postgres

import (
	"context"
	"encoding/json"
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
	var data = &domain.ServicesRequests{
		ID:             uuid.New().String(),
		ProfessionalID: dto.ProfessionalID,
		ClientID:       dto.ClientID,
		Description:    dto.Description,
		Category:       categoryRes,
		LastClientLat:  &clientLat,
		LastClientLng:  &clientLng,
		Status:         statusRes,
	}
	tx := s.db.WithContext(ctx).Model(&domain.ServicesRequests{})
	if dto.ProfessionalID == "" {
		tx = tx.Raw("INSERT INTO services_requests (id, client_id, description, category_id, last_client_lat, last_client_lng, status_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())", data.ID, data.ClientID, data.Description, data.Category.ID, data.LastClientLat, data.LastClientLng, statusRes.ID)
		tx = tx.Select("id, client_id, professional_id, description, last_client_lat as client_latitude, last_client_lng as client_longitude").Scan(&resp)
	} else {
		tx = tx.Create(
			&data,
		).Select("id, client_id, professional_id, description, last_client_lat as client_latitude, last_client_lng as client_longitude").Scan(&resp)
	}

	if tx.Error != nil {
		s.log.Printf("Error inserting service into the database: %v", tx.Error)
		return repositories.CreateServiceResponseDto{}, tx.Error
	}

	resp.CategoryName = categoryRes.Name
	resp.Status = string(statusRes.Name)
	return resp, nil

}

func (s *ServiceRepositoryImp) ValidateExistingPendingServiceFromClientToProfessional(ctx context.Context, clientID string, professionalID string) (bool, error) {
	s.log.Println("Validating existing pending service from client to professional...")
	var count int64
	var statusRes domain.Status
	res := s.db.WithContext(ctx).Model(&domain.Status{}).Where("name = ?", "FINALIZADO").First(&statusRes)
	if res.Error == gorm.ErrRecordNotFound {
		return false, gorm.ErrRecordNotFound
	}
	res = s.db.WithContext(ctx).Model(&domain.ServicesRequests{}).Where("client_id = ? AND professional_id = ? AND status_id != ?", clientID, professionalID, statusRes.ID).Count(&count)

	if res.Error != nil {
		s.log.Printf("Error validating existing pending service from client to professional: %v", res.Error)
		return false, res.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (s *ServiceRepositoryImp) GetClientById(ctx context.Context, clientID string) (*domain.Users, error) {
	var user domain.Users
	res := s.db.WithContext(ctx).Model(&domain.Users{}).Where("id = ?", clientID).First(&user)
	if res.Error != nil {
		s.log.Printf("Error getting client by ID: %v", res.Error)
		return nil, res.Error
	}
	return &user, nil
}

func (s *ServiceRepositoryImp) CountServicesByUserId(ctx context.Context, userID string) (int, error) {
	var count int64
	res := s.db.WithContext(ctx).Model(&domain.ServicesRequests{}).Where("professional_id = ? OR client_id = ?", userID, userID).Count(&count)
	if res.Error != nil {
		s.log.Printf("Error counting services by user ID: %v", res.Error)
		return 0, res.Error
	}
	return int(count), nil
}
func (s *ServiceRepositoryImp) GetServicesByUserId(ctx context.Context, userID string, offset int, limit int) ([]repositories.ServiceDataResponseDto, error) {
	var results []repositories.ServiceDataResponseDto
	var services []domain.ServicesRequests
	res := s.db.WithContext(ctx).Model(&domain.ServicesRequests{}).
		Preload("Professional").Preload("Client").Preload("Status").Preload("Category").
		Where("professional_id = ? OR client_id = ? AND professional_id is not null", userID, userID).
		Offset(offset).Limit(limit).Find(&services)
	if res.Error != nil {
		s.log.Printf("Error getting services by user ID: %v", res.Error)
		return nil, res.Error
	}
	for _, service := range services {
		mappedResult := repositories.ServiceDataResponseDto{
			Status:           string(service.Status.Name),
			Description:      service.Description,
			ProfessionalName: service.Professional.Name + " " + service.Professional.LastName,
			CategoryName:     service.Category.Name,
			ID:               service.ID,
			DateOfCreation:   service.CreatedAt.Format("2006-01-02 15:04:05"),
			Price:            service.AgreedPrice,
		}
		results = append(results, mappedResult)
	}
	return results, nil
}

func (s *ServiceRepositoryImp) GetServiceDetailById(ctx context.Context, serviceID string) (*domain.ServicesRequests, error) {
	var result domain.ServicesRequests

	tx := s.db.WithContext(ctx).Model(&domain.ServicesRequests{})

	tx = tx.Preload("Professional").Preload("Client").Preload("Status").Preload("Category").Preload("ServiceEvidence").Preload("Payments")
	tx = tx.Where("id=?", serviceID).First(&result)
	if tx.Error != nil {
		s.log.Printf("Error getting service detail by ID: %v", tx.Error)
		return nil, tx.Error
	}

	return &result, nil

}

func (s *ServiceRepositoryImp) SaveServiceEvidence(ctx context.Context, dto repositories.SaveServiceEvidenceRequestDto) error {

	tx := s.db.WithContext(ctx).Model(&domain.ServiceEvidence{})
	data, err := json.Marshal(dto.StrokesData)
	if err != nil {
		s.log.Printf("Error marshalling strokes data: %v", err)
		return err
	}
	evidence := &domain.ServiceEvidence{
		ID:          uuid.NewString(),
		ServiceId:   dto.ServiceID,
		SignedByID:  dto.ClientID,
		JsonPayload: data,
	}
	tx.Create(evidence)
	return tx.Error
}

func (s *ServiceRepositoryImp) SaveServiceReview(ctx context.Context, dto repositories.SaveServiceReviewRequestDto) error {
	var service domain.ServicesRequests
	res := s.db.WithContext(ctx).Model(&domain.ServicesRequests{}).
		Where("id = ? AND client_id = ?", dto.ServiceId, dto.ClientId).First(&service)

	if res.Error != nil {
		s.log.Printf("Error getting service by ID and client ID: %v", res.Error)
		return res.Error
	}

	review := s.db.WithContext(ctx).Model(&domain.ServiceReviews{}).
		Create(&domain.ServiceReviews{
			ID:                uuid.NewString(),
			ServiceRequestsId: dto.ServiceId,
			ReviewerId:        dto.ClientId,
			RevieweedId:       service.ProfessionalID,
			Rating:            dto.Rating,
			Comment:           dto.Comment,
		})

	if review.Error != nil {
		s.log.Printf("Error saving service review: %v", review.Error)
		return review.Error
	}
	return nil

}
