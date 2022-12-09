package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"control-pago-backend/internal/email"
	"control-pago-backend/internal/repository"
	"control-pago-backend/internal/service"
	"control-pago-backend/internal/transport"
	"control-pago-backend/log"
)

var (
	// Api
	ApiPort = os.Getenv("API_PORT")
	// Postgres
	PostgresUser        = os.Getenv("POSTGRES_USER")
	PostgresPwd         = os.Getenv("POSTGRES_PWD")
	PostgresHost        = os.Getenv("POSTGRES_HOST")
	PostgresPort        = os.Getenv("POSTGRES_PORT")
	PostgresDB          = os.Getenv("POSTGRES_DB")
	EmailPwd            = os.Getenv("EMAIL_PWD")
	ReceiptsFolderPath  = os.Getenv("RECEIPTS_FOLDER_PATH")
	RegisterPaymentpath = os.Getenv("REGISTER_PAYMENT_PATH")
	GetSummaryPath      = os.Getenv("GET_SUMMARY_PATH")
	StaticServerBaseUrl = os.Getenv("STATIC_SERVER_BASE_URL")
)

func main() {
	logger := log.NewLogger()

	// create postgres connection
	postgresConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", PostgresUser, PostgresPwd, PostgresHost, PostgresPort, PostgresDB)
	postgresDb, err := sql.Open("postgres", postgresConnString)
	if err != nil {
		logger.Error("main.go", "main", err.Error())
		panic(err)
	}
	defer postgresDb.Close()

	err = postgresDb.Ping()
	if err != nil {
		logger.Error("main.go", "main", err.Error())
		panic(err)
	}

	repository := repository.NewRepository(logger, postgresDb, ReceiptsFolderPath)

	emailClient := email.NewEmailClient(EmailPwd)

	service := service.NewService(emailClient, logger, repository, StaticServerBaseUrl)

	handlers := transport.NewHandler(service, logger)

	// start mux and listening
	router := transport.NewRouter(handlers, RegisterPaymentpath, GetSummaryPath, logger)
	addr := fmt.Sprintf("0.0.0.0:%s", ApiPort)
	logger.Info("main.go", "main", "Started listening on "+addr)
	go http.ListenAndServe(addr, router)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill, syscall.SIGKILL)
	<-c

	logger.Info("main.go", "main", "program interrupted, gracefully shuting down")
	os.Exit(1)
}
