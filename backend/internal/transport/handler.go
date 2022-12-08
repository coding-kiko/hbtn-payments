package transport

import (
	"control-pago-backend/internal/entity"
	"control-pago-backend/internal/errors"
	"control-pago-backend/internal/service"
	"control-pago-backend/log"
	"encoding/json"
	"net/http"
)

type handler struct {
	service service.Service
	logger  log.Logger
}

type Handler interface {
	RegisterPayment(w http.ResponseWriter, r *http.Request)
	// GetSummary(w http.ResponseWriter, r *http.Request)

	MethodNotAllowedHandler() http.Handler
}

func NewHandler(srv service.Service, lgr log.Logger) Handler {
	return &handler{
		service: srv,
		logger:  lgr,
	}
}

func (h *handler) RegisterPayment(w http.ResponseWriter, r *http.Request) {
	req := new(entity.RegisterPaymentRequest)

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		h.logger.Error("handler.go", "RegisterPayment", err.Error())
		statusCode, body := errors.CreateResponse(errors.NewBadRequest("bad request: decoding body"))
		makeHttpRensponse(w, statusCode, body)
		return
	}

	err = h.service.RegisterPayment(req)
	if err != nil {
		h.logger.Error("handler.go", "RegisterPayment", err.Error())
		statusCode, body := errors.CreateResponse(err)
		makeHttpRensponse(w, statusCode, body)
		return
	}
	w.WriteHeader(204)
}

// override default gorilla method not allowed handler
func (h *handler) MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		statusCode, body := errors.CreateResponse(errors.NewMethodNotAllowed())
		makeHttpRensponse(w, statusCode, body)
	})
}

func makeHttpRensponse(w http.ResponseWriter, statusCode int, body errors.Response) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
