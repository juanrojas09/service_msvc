package repositories

import (
	"context"

	"github.com/juanrojas09/core_domain/domain"
)

type CreateServiceRequestDTO struct {
	ProfessionalID  string `json:"professional_id" `
	ClientID        string `json:"client_id" validate:"required"`
	Description     string `json:"description" validate:"required"`
	CategoryID      string `json:"category_id"`
	ClientLatitude  string `json:"client_latitude"`
	ClientLongitude string `json:"client_longitude"`
}

type CreateServiceResponseDto struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id"`
	ProfessionalID  string `json:"professional_id"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	CategoryName    string `json:"category_name"`
	ClientLatitude  string `json:"last_client_lat"`
	ClientLongitude string `json:"last_client_lng"`
}

type ServiceRepository interface {
	// Define the methods that the service repository should have
	CreateService(ctx context.Context, dto CreateServiceRequestDTO) (CreateServiceResponseDto, error)
	GetClientById(ctx context.Context, clientID string) (*domain.Users, error)

	ValidateExistingPendingServiceFromClientToProfessional(ctx context.Context, clientID string, professionalID string) (bool, error)
}
