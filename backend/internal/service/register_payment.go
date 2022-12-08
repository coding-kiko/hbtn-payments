package service

import (
	"crypto/md5"
	enc "encoding/base64"
	"fmt"
	"time"

	"control-pago-backend/internal/entity"
	"control-pago-backend/internal/errors"
)

func (s *service) RegisterPayment(req *entity.RegisterPaymentRequest) error {
	var receiptBase64 = req.Receipt

	// Decode receipt
	content, err := enc.StdEncoding.DecodeString(receiptBase64)
	if err != nil {
		s.logger.Error("register_payment.go", "RegisterPayment", err.Error())
		return errors.NewFileError("Error decoding base64 file string")
	}

	// generate link name
	// ext, err := filetype.Match(content)
	// if err != nil {
	// 	s.logger.Error("register_payment.go", "RegisterPayment", err.Error())
	// 	return errors.NewFileError("Error getting file extension")
	// }
	fileName := generateReceiptFileName(receiptBase64[:8], req.Month, "jpg")

	receipt := entity.Receipt{
		Name: fileName,
		Data: content,
	}

	// store receipt disk
	err = s.Repo.StoreReceipt(receipt)
	if err != nil {
		s.logger.Error("register_payment.go", "RegisterPayment", err.Error())
		return err
	}

	req.Receipt = fileName
	// register payment psql
	err = s.Repo.RegisterPayment(req)
	if err != nil {
		s.logger.Error("register_payment.go", "RegisterPayment", err.Error())
		return err
	}

	err = s.EmailClient.SendReceipt(req.EmailTo, receiptBase64, req.Month)
	if err != nil {
		s.logger.Error("register_payment.go", "RegisterPayment", err.Error())
		return err
	}

	return nil
}

// concatenate current time, an arbitrary chunk of the base64 encrypted file and the month
// then hash it with md5 to produce a random string that will be the file name
func generateReceiptFileName(receiptChunk, month, extension string) string {
	now := time.Now().String()

	b := []byte(now + receiptChunk + month)
	filename := fmt.Sprintf("%x", md5.Sum(b))
	return fmt.Sprintf("%s.%s", filename, extension)
}
