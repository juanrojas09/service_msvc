package usecases

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"strconv"

	"github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
)

type (
	ServiceDetailByIdImpl struct {
		repo repositories.ServiceRepository
		log  *log.Logger
	}
)

func NewServiceDetailByIdImpl(repo repositories.ServiceRepository, logger *log.Logger) interfaces.UseCases {
	return &ServiceDetailByIdImpl{
		repo: repo,
		log:  logger,
	}
}

func (s *ServiceDetailByIdImpl) Handle(ctx context.Context, params ...interface{}) (interface{}, error) {
	s.log.Println("Getting service details by ID...")
	req := params[0].(repositories.ServiceDetailRequestDto)
	var mappedResult repositories.ServiceDetailResponseDto

	if req.ID == "" {
		return nil, response.BadRequest("El ID del servicio es obligatorio")
	}

	result, err := s.repo.GetServiceDetailById(ctx, req.ID)
	if err != nil {
		return nil, response.InternalServerError(err.Error())
	}

	mappedResult.ID = result.ID
	mappedResult.ServiceDetail = repositories.ServiceDataDetail{
		CategoryName:     result.Category.Name,
		Description:      result.Description,
		Status:           string(result.Status.Name),
		ProfessionalName: result.Professional.Name + " " + result.Professional.LastName,
		DateOfCreation:   result.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	mappedResult.IsMatched = result.ProfessionalID != "" && (result.Status.Name == "CONFIRMADO" || result.Status.Name == "EN PROGRESO" || result.Status.Name == "PENDIENTE DE PAGO" || result.Status.Name == "FINALIZADO")
	mappedResult.IsBudgetOffered = result.AgreedPrice > 0

	//map payments
	for _, payment := range result.Payments {
		payment := repositories.PaymentData{
			Amount: payment.Amount,
			Status: string(payment.PaymentStatus.Name),
			Date:   payment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		mappedResult.PaymentData = append(mappedResult.PaymentData, payment)
	}

	//map timeline
	mappedResult.TimeLineData = append(mappedResult.TimeLineData, repositories.TimeLineData{
		Title:   "Solicitud de servicio creada",
		Message: "El cliente " + result.Client.Name + " " + result.Client.LastName + " ha creado una solicitud de servicio.",
	})

	if result.ProfessionalID != "" && (result.Status.Name == "CONFIRMADO" || result.Status.Name == "EN PROGRESO" || result.Status.Name == "PENDIENTE DE PAGO" || result.Status.Name == "FINALIZADO") {
		mappedResult.TimeLineData = append(mappedResult.TimeLineData, repositories.TimeLineData{
			Title:   "Servicio Confirmado",
			Message: "El profesional " + result.Professional.Name + " " + result.Professional.LastName + " ha confirmado el servicio.",
		})
	}

	if result.AgreedPrice > 0 && result.AgreedPriceAt != nil {
		mappedResult.TimeLineData = append(mappedResult.TimeLineData, repositories.TimeLineData{
			Title:   "Presupuesto acordado",
			Message: "El presupuesto acordado es de $" + strconv.FormatFloat(math.Round(result.AgreedPrice), 'f', 2, 64) + ".",
		})
	}

	if result.ProfessionalID != "" && (result.Status.Name == "EN PROGRESO" || result.Status.Name == "PENDIENTE DE PAGO" || result.Status.Name == "FINALIZADO") {
		mappedResult.TimeLineData = append(mappedResult.TimeLineData, repositories.TimeLineData{
			Title:   "Profesional En Camino",
			Message: "El profesional " + result.Professional.Name + " " + result.Professional.LastName + " está en camino.",
		})
	}

	if result.ProfessionalID != "" && (result.Status.Name == "PENDIENTE DE PAGO" || result.Status.Name == "FINALIZADO") {
		mappedResult.TimeLineData = append(mappedResult.TimeLineData, repositories.TimeLineData{
			Title:   "Servicio En Progreso",
			Message: "El profesional " + result.Professional.Name + " " + result.Professional.LastName + " ha iniciado el servicio.",
		})
	}

	if result.Status.Name == "FINALIZADO" {
		mappedResult.TimeLineData = append(mappedResult.TimeLineData, repositories.TimeLineData{
			Title:   "Servicio Finalizado",
			Message: "El servicio ha sido finalizado.",
		})
	}

	//Mapear evidencias
	if result.ServiceEvidence.ID != "" {
		if err := json.Unmarshal(result.ServiceEvidence.JsonPayload, &mappedResult.ServiceEvidence); err != nil {
			return nil, response.InternalServerError(err.Error())
		}
	}

	return response.Ok("Service details fetched successfully", mappedResult, nil), nil
}
