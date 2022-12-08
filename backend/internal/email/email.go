package email

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/smtp"

	"control-pago-backend/internal/errors"
)

type EmailClient interface {
	SendReceipt(to, receiptBase64, month string) error
}

type emailClient struct {
	SmtpAddr string
	Auth     *smtp.Auth
	From     string
}

func NewEmailClient(pwd string) EmailClient {
	var smtpHost = "smtp.gmail.com"
	var smtpPort = "587"
	var from = "kikoaudi2001@gmail.com"

	auth := smtp.PlainAuth("", from, pwd, smtpHost)

	return &emailClient{
		Auth:     &auth,
		SmtpAddr: fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		From:     from,
	}
}

func (e *emailClient) SendReceipt(to, receiptBase64, month string) error {
	content := makeEmailContent(receiptBase64, month)
	err := smtp.SendMail(e.SmtpAddr, *e.Auth, e.From, []string{to}, content)
	if err != nil {
		return errors.NewInternalServer(fmt.Sprintf("Error sending email to %s", to))
	}
	return nil
}

func makeEmailContent(receiptBase64, month string) []byte {
	var buf = bytes.NewBuffer(nil)

	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	mimeMsgBoundary := fmt.Sprintf("\n--%s\n", boundary)

	buf.WriteString(fmt.Sprintf("Subject: Comprobante %s\n", month))
	buf.WriteString(fmt.Sprintf("MIME-version: 1.0;\nContent-Type: multipart/mixed; boundary=%s", boundary))
	buf.WriteString(mimeMsgBoundary)
	buf.WriteString(fmt.Sprintf("Buen dia,\nAdjunto el comprobante del %s.\n\nSaludos,\nFrancisco Calixto.", month))
	buf.WriteString(mimeMsgBoundary)
	buf.WriteString("Content-Type: image/jpg\n")
	buf.WriteString("Content-Transfer-Encoding: base64\n")
	buf.WriteString("Content-Disposition: attachment; filename=comprobante\n")
	buf.WriteString(receiptBase64)

	return buf.Bytes()
}
