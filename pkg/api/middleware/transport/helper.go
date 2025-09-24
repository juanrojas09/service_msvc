package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/service_msvc/pkg/api/app/controllers"
	"github.com/juanrojas09/service_msvc/pkg/api/app/interfaces/repositories"
)

func NewHttpServer(ctx context.Context, endpoints controllers.Endpoints) http.Handler {
	r := gin.Default()

	CreateServiceRequestHandler := gin.WrapH(kithttp.NewServer(endpoint.Endpoint(endpoints.CreateServiceRequest),
		decodeCreateServiceRequest, encodeResponse))

	ListServiceRequestHandler := gin.WrapH(kithttp.NewServer(endpoint.Endpoint(endpoints.GetServiceListByUserIdRequest),
		decodeListServiceByUserIdRequest, encodeResponse))

	r.POST("/service", encodeGin, CreateServiceRequestHandler)
	r.GET("/service/:id", encodeGin, ListServiceRequestHandler)
	return r
}

func encodeGin(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "params", c.Params)
	c.Request = c.Request.WithContext(ctx)
}
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateServiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := repositories.CreateServiceRequestDTO{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func decodeListServiceByUserIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := repositories.ServiceListRequestDTO{}

	params := ctx.Value("params").(gin.Params)

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	req.UserID = params.ByName("id")

	if limit != "" {
		req.Limit, _ = strconv.Atoi(limit)
	}

	if offset != "" {
		req.Page, _ = strconv.Atoi(offset)
	}

	return req, nil
}
