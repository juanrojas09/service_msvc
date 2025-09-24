package usecases

import (
	"context"
	"encoding/json"
	"log"

	"github.com/juanrojas09/core_domain/domain"
	"github.com/juanrojas09/core_rabbit_pub_sub_base/broker"
	r "github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
	"github.com/rabbitmq/amqp091-go"
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
		s.log.Println("Ya exoste una solicitud de servicio pendiente entre el cliente y el profesional.")
		return repositories.CreateServiceResponseDto{}, r.BadRequest("Ya exoste una solicitud de servicio pendiente entre el cliente y el profesional.")
	}
	response, err := s.repo.CreateService(ctx, req)
	if err != nil {
		s.log.Printf("Error creating service: %v", err)
		return repositories.CreateServiceResponseDto{}, r.InternalServerError(err.Error())
	}

	var config broker.RabbitMQConfig
	rabbitConfig := config.New("localhost", "kalo", "kalo", "5672")
	queue, rabbit, err := s.ConnectToRabbit(rabbitConfig)
	if err != nil {
		s.log.Printf("Error connecting to RabbitMQ: %v", err)
		return repositories.CreateServiceResponseDto{}, r.InternalServerError(err.Error())
	}

	client, err := s.repo.GetClientById(ctx, req.ClientID)
	if err != nil {
		s.log.Printf("Error getting client by ID: %v", err)
		return repositories.CreateServiceResponseDto{}, r.InternalServerError(err.Error())
	}

	defer func() {

		var notification = domain.Notifications{
			ID:      "",
			UserID:  req.ProfessionalID,
			Title:   "Nuevo Servicio",
			Message: "Tienes una nueva peticion de servicio del cliente: " + client.Name + ". Descripcion: " + req.Description,
		}
		messageBody, err := json.Marshal(&notification)
		if err != nil {
			s.log.Printf("Error marshalling event: %v", err)
		}
		err = rabbit.Publish(queue.Name, string(messageBody))
		if err != nil {
			s.log.Printf("Error publishing event: %v", err)
		}
	}()

	return r.Created("Service Created Successfully", response, nil), nil
}

func (s *ServiceRequestImpl) ConnectToRabbit(config *broker.RabbitMQConfig) (amqp091.Queue, *broker.RabbitMQ, error) {
	rabbit := broker.NewRabbitMQ(*config)
	err := rabbit.Connect()
	if err != nil {
		s.log.Printf("Error connecting to RabbitMQ: %v", err)
		return amqp091.Queue{}, nil, err
	}

	err = rabbit.Channel()
	if err != nil {
		s.log.Printf("Error creating channel to RabbitMQ: %v", err)
		return amqp091.Queue{}, nil, err
	}
	queue, err := rabbit.DeclareQueue("NewNotificationEvent")
	if err != nil {
		s.log.Printf("Error declaring queue in RabbitMQ: %v", err)
		return amqp091.Queue{}, nil, err
	}
	return queue, rabbit, nil
}
