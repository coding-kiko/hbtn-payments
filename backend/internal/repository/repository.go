package repository

import (
	"database/sql"
	"fmt"
	"os"

	"control-pago-backend/internal/entity"
	"control-pago-backend/internal/errors"
	"control-pago-backend/log"
)

type Repository interface {
	RegisterPayment(pmt *entity.RegisterPayment) error
	StoreReceipt(receipt entity.Receipt) error
	// GetPayments()
}

var (
	registerPaymentQuery = `INSERT INTO payments(month, amount, receipt_url, company)
										VALUES($1, $2, $3, $4)`
)

type repository struct {
	logger         log.Logger
	db             *sql.DB
	receiptsFolder string
}

func (r *repository) RegisterPayment(pmt *entity.RegisterPayment) error {

	_, err := r.db.Exec(registerPaymentQuery, pmt.Month, pmt.Amount, pmt.Receipt, pmt.Company)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) StoreReceipt(receipt entity.Receipt) error {
	f, err := os.Create(fmt.Sprintf("%s/%s", r.receiptsFolder, receipt.Name))
	if err != nil {
		r.logger.Error("repository.go", "StoreReceipt", err.Error())
		return errors.NewFileError("Error creating file")
	}
	defer f.Close()

	_, err = f.Write(receipt.Data)
	if err != nil {
		r.logger.Error("repository.go", "StoreReceipt", err.Error())
		return errors.NewFileError("Error writing content to file")
	}

	return nil
}

func NewRepository(lgr log.Logger, db *sql.DB, receiptsFolder string) Repository {
	return &repository{
		logger:         lgr,
		db:             db,
		receiptsFolder: receiptsFolder,
	}
}
