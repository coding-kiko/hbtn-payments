package transport

import (
	"control-pago-backend/log"

	"github.com/gorilla/mux"
)

func NewRouter(handler Handler, registerPaymentPath, getSummaryPath string, logger log.Logger) *mux.Router {
	router := mux.NewRouter()
	logger.Info("router.go", "NewRouter", "Initializing handlers")

	router.Path(registerPaymentPath).Methods("POST").HandlerFunc(handler.RegisterPayment)
	router.Path(getSummaryPath).Methods("GET").HandlerFunc(handler.GetSummary)

	// override default gorilla 405 handler
	router.MethodNotAllowedHandler = handler.MethodNotAllowedHandler()

	return router
}
