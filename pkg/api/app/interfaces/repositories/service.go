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

type ServiceListRequestDTO struct {
	UserID string `json:"user_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type ServiceListResponseDTO struct {
	Data []ServiceDataResponseDto `json:"data"`
}

type ServiceDataResponseDto struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	Description      string  `json:"description"`
	ProfessionalName string  `json:"professional_name"`
	CategoryName     string  `json:"category_name"`
	DateOfCreation   string  `json:"created_at"`
	Price            float64 `json:"price"`
}

type ServiceDetailRequestDto struct {
	ID string `json:"id"`
}

type ServiceDetailResponseDto struct {
	ID              string            `json:"id"`
	ServiceDetail   ServiceDataDetail `json:"service_detail"`
	IsMatched       bool              `json:"is_matched"`
	IsBudgetOffered bool              `json:"is_budget_offered"`
	TimeLineData    []TimeLineData    `json:"time_line_data"`
	PaymentData     []PaymentData     `json:"payment_data"`
	ServiceEvidence EvidenceJSON      `json:"service_evidence"`
}

type ServiceDataDetail struct {
	CategoryName     string `json:"category_name"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	ProfessionalName string `json:"professional_name"`
	DateOfCreation   string `json:"date_of_creation"`
}

type TimeLineData struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PaymentData struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
	Date   string  `json:"date"`
}
type EvidenceJSON struct {
	Canvas struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"canvas"`
	PenColor string `json:"penColor"`
	Strokes  []struct {
		Thickness float64     `json:"thickness"`
		Points    [][]float64 `json:"points"`
	} `json:"strokes"`
}
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	T int     `json:"t"`
}

type SaveServiceEvidenceRequestDto struct {
	ServiceID   string       `json:"service_id"`
	ClientID    string       `json:"client_id"`
	StrokesData EvidenceJSON `json:"strokes_data"`
}

type ServiceRepository interface {
	// Define the methods that the service repository should have
	CreateService(ctx context.Context, dto CreateServiceRequestDTO) (CreateServiceResponseDto, error)
	GetClientById(ctx context.Context, clientID string) (*domain.Users, error)

	ValidateExistingPendingServiceFromClientToProfessional(ctx context.Context, clientID string, professionalID string) (bool, error)
	CountServicesByUserId(ctx context.Context, userID string) (int, error)
	GetServicesByUserId(ctx context.Context, userID string, offset int, limit int) ([]ServiceDataResponseDto, error)
	GetServiceDetailById(ctx context.Context, serviceID string) (*domain.ServicesRequests, error)

	SaveServiceEvidence(ctx context.Context, dto SaveServiceEvidenceRequestDto) error
}
