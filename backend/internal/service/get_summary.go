package service

import (
	"bytes"
	"control-pago-backend/internal/entity"
	enc "encoding/base64"
	"fmt"
)

func (s *service) GetSummary() (*entity.GetSummaryResponse, error) {
	payments, err := s.Repo.GetPayments()
	if err != nil {
		s.logger.Error("get_summay.go", "GetSummary", err.Error())
		return nil, err
	}

	summaryBase64, err := GenerateHtmlSummaryBase64(payments, s.StaticServerBaseUrl)
	if err != nil {
		s.logger.Error("get_summay.go", "GetSummary", err.Error())
		return nil, err
	}

	res := &entity.GetSummaryResponse{
		Summary: summaryBase64,
	}

	return res, nil
}

func GenerateHtmlSummaryBase64(payments []entity.Payment, staticServerBaseUrl string) (string, error) {
	var total int = 0
	var file = bytes.NewBuffer(nil)

	file.WriteString(fixedTop)
	for _, payment := range payments {
		total += payment.Amount
		row := fmt.Sprintf(dynamicRow, payment.Month, payment.Amount, payment.Company, staticServerBaseUrl+payment.Receipt)
		file.WriteString(row)
	}
	file.WriteString(fmt.Sprintf(fixedBot, total))

	fileBase64 := enc.StdEncoding.EncodeToString(file.Bytes())

	return fileBase64, nil
}

var fixedTop = `<!DOCTYPE html>
				<html>
				<head>
					<style>
						.styled-table {
							border-collapse: collapse;
							margin: 25px 0;
							font-size: 0.9em;
							font-family: sans-serif;
							min-width: 400px;
							box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
						}
						.styled-table thead tr {
							background-color: #353d3c;
							color: #ffffff;
							text-align: left;
						}
						.styled-table th,
						.styled-table td {
							padding: 12px 15px;
						}
						.styled-table tbody tr {
							border-bottom: 1px solid #dddddd;
						}
						.styled-table tbody tr:nth-of-type(even) {
							background-color: #f3f3f3;
						}
						.styled-table tbody tr:last-of-type {
							border-bottom: 2px solid #353d3c;
						}
					</style>
				</head>
				<body>
					<table class="styled-table">
						<thead>
							<tr class="active-row">
								<th>Month</th>
								<th>Amount</th>
								<th>Company</th>
								<th>Receipt</th>
							</tr>
						</thead>
						<tbody>`

var fixedBot = `		</tbody>
					</table>
					<h2>Total Payed: %d</h2>
				</body>

				</html>`

var dynamicRow = `<tr>
					<td>%s</td>
					<td>%d</td>
					<td>%s</td>
					<td><a href="%s">view</td>
				   </tr>`
