package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/juanrojas09/core_domain/common"
	c "github.com/juanrojas09/core_sdk/auth"
	"github.com/juanrojas09/service_msvc/pkg/api/app/controllers"
	"github.com/juanrojas09/service_msvc/pkg/api/app/usecases"
	"github.com/juanrojas09/service_msvc/pkg/api/configs/bootstrap"
	"github.com/juanrojas09/service_msvc/pkg/api/middleware/transport"
	"github.com/juanrojas09/service_msvc/pkg/persistance/postgres"
)

func main() {
	godotenv.Load()

	db := bootstrap.InitDatabase()
	logger := bootstrap.InitLogger()

	//dependency container
	serviceRepository := postgres.NewServiceRepository(db, logger)
	createServiceUC := usecases.NewServiceRequestImpl(serviceRepository, logger)
	listServicesUC := usecases.NewServiceListByUserIdImpl(serviceRepository, logger)
	GetServiceDetailByIdUC := usecases.NewServiceDetailByIdImpl(serviceRepository, logger)
	SaveServiceEvidenceUC := usecases.NewSaveServiceEvidenceImpl(serviceRepository, logger)
	SaveServiceReviewsUC := usecases.NewSaveServiceReviewsImpl(serviceRepository, logger)
	jwtService := common.NewJWTService()

	registry := &controllers.UseCaseRegistry{
		CreateServiceRequestUseCase: createServiceUC,
		ListServiceByUserIdUseCase:  listServicesUC,
		GetServiceDetailByIdUseCase: GetServiceDetailByIdUC,
		SaveServiceEvidenceUseCase:  SaveServiceEvidenceUC,
		SaveServiceReviewsUseCase:   SaveServiceReviewsUC,
	}
	endpoints := controllers.MakeEndpoints(registry)

	ctx := context.Background
	h := transport.NewHttpServer(ctx(), endpoints)
	Addr := os.Getenv("API_URL_BASE") + ":" + os.Getenv("API_PORT")

	srv := http.Server{
		Handler: setUpBaseHeaders(h, jwtService),
		Addr:    Addr,
	}

	chnErr := make(chan error)
	go func() {
		log.Println("Server listening on port", Addr)
		chnErr <- srv.ListenAndServe()

	}()
	err := <-chnErr
	if err != nil {
		log.Fatal(err)
	}

}

func setUpBaseHeaders(h http.Handler, jwtService common.JWTService) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if strings.Contains(r.URL.Path, "auth") {
			h.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Authorization header"})
			return
		}

		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		headers.Set("Accept", "application/json")
		headers.Set("Authorization", "Bearer "+tokenString)
		transport := c.NewClientHttp("http://localhost:8000", headers)
		res, error := transport.Get("/auth/trust")
		if error != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "" + error.Error()})
			return
		}

		if res.Code != http.StatusOK {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
			return
		}

		h.ServeHTTP(w, r)

	})
}
