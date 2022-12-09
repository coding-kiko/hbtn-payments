package bot

import (
	"bytes"
	"context"
	encoding "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn               *whatsmeow.Client
	groupJId           types.JID
	stage              int
	RequestBody        ReqBody
	getSummaryUri      string
	registerPaymentUri string
}

func NewClient(conn *whatsmeow.Client, groupJid, registerPaymentUri, getSummaryUri string) *Client {
	jid, _ := types.ParseJID(groupJid)
	return &Client{
		conn:               conn,
		stage:              0,
		groupJId:           jid,
		registerPaymentUri: registerPaymentUri,
		getSummaryUri:      getSummaryUri,
	}
}

func (c *Client) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Message.GetConversation() == "get-summary" && v.Info.Chat == c.groupJId && c.stage == 0 {
			c.SendMessage("<BOT> Getting payments summary...")
			summaryBase64, err := c.GetPaymentsSummaryRequest()
			if err != nil {
				c.SendMessage("<BOT> Could not get payments summary")
			}

			summary, _ := encoding.StdEncoding.DecodeString(summaryBase64)
			err = c.SendHtmlFile(summary)
			if err != nil {
				c.SendMessage(fmt.Sprintf("<BOT> Error getting summary: %s", err.Error()))
			}

			return
		}
		if v.Message.GetConversation() == "register-payment" && v.Info.Chat == c.groupJId && c.stage == 0 {
			c.SendMessage("<BOT> Indicate amount payed [integer]:")
			c.stage = 1
			return
		}
		if c.stage == 1 && v.Info.Chat == c.groupJId {
			msg := v.Message.GetConversation()
			if messageIsFromBot(msg) {
				return
			}

			amount, err := strconv.Atoi(msg)
			if err != nil {
				c.SendMessage("<BOT> Please enter a numeric value")
				return
			}
			c.RequestBody.Amount = amount

			c.SendMessage("<BOT> Indicate month [MM/YYYY]:")
			c.stage = 2
			return
		}
		if c.stage == 2 && v.Info.Chat == c.groupJId {
			month := v.Message.GetConversation()
			if messageIsFromBot(month) {
				return
			}

			match, _ := regexp.MatchString(`[0-1][0-9]/[0-9]{4}`, month)
			if !match {
				c.SendMessage("<BOT> Please match the format 'MM/YYYY'")
				return
			}
			c.RequestBody.Month = month

			c.SendMessage("<BOT> Indicate email receiver [default=mvd-accounting@holbertonschool.com | none='dont send email']:")
			c.stage = 3
			return
		}
		if c.stage == 3 && v.Info.Chat == c.groupJId {
			var defaultEmail = "mvd-accounting@holbertonschool.com"

			email := v.Message.GetConversation()
			if messageIsFromBot(email) {
				return
			}

			switch {
			case email == "default":
				c.RequestBody.EmailTo = &defaultEmail
			case email == "none":
				c.RequestBody.EmailTo = nil
			default:
				c.RequestBody.EmailTo = &email
			}

			c.SendMessage("<BOT> Indicate company [none=undisclosed]:")
			c.stage = 4
			return
		}
		if c.stage == 4 && v.Info.Chat == c.groupJId {
			company := v.Message.GetConversation()
			if messageIsFromBot(company) {
				return
			}

			if company == "none" {
				c.RequestBody.Company = "N/A"
			} else {
				c.RequestBody.Company = company
			}

			c.SendMessage("<BOT> Insert image of the receipt:")
			c.stage = 5
			return
		}
		if c.stage == 5 && v.Info.Chat == c.groupJId {
			if v.Message.ImageMessage == nil {
				c.SendMessage("<BOT> Error not an image")
				return
			}

			data, err := c.conn.Download(v.Message.GetImageMessage())
			if err != nil {
				c.SendMessage("<BOT> Error downloading receipt")
				c.stage = 0
			}

			receipt := encoding.StdEncoding.EncodeToString(data)
			c.RequestBody.ReceiptBASE64 = receipt

			err = c.PostRegisterPaymentRequest()
			if err != nil {
				c.SendMessage(fmt.Sprintf("<BOT> Error storing the payment: %s", err.Error()))
				c.stage = 0
				return
			}

			if c.RequestBody.EmailTo == nil {
				c.SendMessage("<BOT> The payment has been registered.")
			} else {
				c.SendMessage("<BOT> The payment has been registered and an email has been sent.")
			}
			c.stage = 0
			return
		}
	}
}

func (c *Client) PostRegisterPaymentRequest() error {
	var buf = bytes.NewBuffer(nil)

	err := json.NewEncoder(buf).Encode(c.RequestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.registerPaymentUri, buf)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return errors.New(fmt.Sprintf("Status Code: %d", res.StatusCode))
	}
	return nil
}

func (c *Client) GetPaymentsSummaryRequest() (string, error) {

	req, err := http.NewRequest("GET", c.getSummaryUri, nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Status Code: %d", res.StatusCode))
	}

	body := GetPaymentSummaryResponse{}

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return "", err
	}
	return body.SummaryBase64, nil
}

// receives summary file as []byte and sneds to the group
func (c *Client) SendHtmlFile(content []byte) error {
	// upload to wpp servers
	resp, err := c.conn.Upload(context.Background(), content, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}

	// generate mediaDocument
	documentMsg := &waProto.DocumentMessage{
		Mimetype:      proto.String("text/html"),
		FileName:      proto.String("summary.html"),
		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}
	_, err = c.conn.SendMessage(context.Background(), c.groupJId, "", &waProto.Message{
		DocumentMessage: documentMsg,
	})
	if err != nil {
		return err
	}

	return nil
}

// syntatic sugar
func (c *Client) SendMessage(msg string) {
	c.conn.SendMessage(context.Background(), c.groupJId, "", &waProto.Message{
		Conversation: proto.String(msg),
	})
}

func messageIsFromBot(msg string) bool {
	return strings.Contains(msg, "<BOT>") || msg == ""
}
