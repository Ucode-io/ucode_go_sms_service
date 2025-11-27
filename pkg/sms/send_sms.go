package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"ucode/ucode_go_sms_service/genproto/sms_service"
	"ucode/ucode_go_sms_service/pkg/models"
	"ucode/ucode_go_sms_service/storage"

	"github.com/google/uuid"
	_ "github.com/lib/pq" //db driver

	"ucode/ucode_go_sms_service/config"

	"github.com/pkg/errors"
)

// Daemon ...
type Daemon struct {
	Conf config.Config
	Strg storage.StorageI
}

const (
	// server we are authorized to send email through
	host                   = "smtp.gmail.com"
	hostPort               = ":587"
	from            string = "ucode.udevs.io@gmail.com"
	defaultPassword string = "xkiaqodjfuielsug"
)

// Init initializes Deamon
func (dmn *Daemon) Init() {

	c := make(chan string)

	isShuttingDown := false
	started := false

OUTTER:
	for {
		if isShuttingDown {
			break
		}
		select {
		case s, ok := <-c:

			if ok {
				fmt.Printf("\nservice is going down...\n")
				fmt.Printf("\nReceived signal %x\n", s)
				isShuttingDown = true
				c = nil
				continue OUTTER
			}
		default:
			var smss []*sms_service.Sms
			smss, err := dmn.Strg.Sms().GetNotSent(context.Background())
			if err != nil {
				log.Println("Error while checking status: ", err)
			}
			for _, sms := range smss {
				fmt.Println("SMS TYPE::::::", sms.Type)
				if strings.ToUpper(sms.Type) == "PHONE" {
					sms.Text = strings.ReplaceAll(sms.Text, "[code]", sms.Otp)
					if strings.Contains(sms.Text, sms.Otp) {
						err = dmn.sendSms(sms.Text, sms.Recipient, sms.DevEmail, sms.DevEmailPassword, sms.Originator)
					} else {
						err = dmn.sendSms(sms.Text+": "+sms.Otp, sms.Recipient, sms.DevEmail, sms.DevEmailPassword, sms.Originator)
					}
				}
				if strings.ToUpper(sms.Type) == "EMAIL" {
					err = dmn.sendCodeToEmail(sms.Text, sms.Recipient, sms.DevEmail, sms.DevEmailPassword)
				}
				if strings.ToUpper(sms.Type) == "MAILCHIMP" {
					err = dmn.sendWithMailchimp(sms.Text, sms.Recipient, sms.DevEmail, sms.DevEmailPassword)
				}
				if err != nil {
					log.Println("error while sending sms: ", err)

					err := dmn.Strg.Sms().IncrementSendCount(context.Background(), sms.Id)
					log.Println("error while incrementing sms send trying count ", err)
				} else {
					err = dmn.Strg.Sms().MakeSent(context.Background(), sms.Id)
					if err != nil {
						log.Println("error while updating status")
					}
				}
			}

			if !started {
				fmt.Print("\n\nSuccessfully\n\n")
				started = true
			}
			t := time.Now()
			fmt.Println(fmt.Sprint("Service is running...\n", t.Format(time.RFC3339)))
			time.Sleep(10 * time.Second)
		}
	}
	fmt.Print("\n\nThis service has finally shutted down\n\n")
}

// Send Mail...
func (dmn *Daemon) sendSms(text, phoneNumber, login, password, originator string) error {
	if len(phoneNumber) < 2 {
		return nil
	}
	phone := phoneNumber[1:]
	return sendWithPlayMobile(dmn, text, phone, login, password, originator)
}

func (dmn *Daemon) sendCodeToEmail(text, to, email, password string) error {
	log.Printf("---SendCodeEmail---> email: %s, code: %s", to, text)

	if email == "" {
		email = from
	}
	if password == "" {
		password = defaultPassword
	}
	message := `
	` + text

	auth := smtp.PlainAuth("", email, password, host)

	msg := "To: \"" + to + "\" <" + to + ">\n" +
		"From: \"" + email + "\" <" + email + ">\n" +
		"Subject: " + "Your verification code" + "\n" +
		message + "\n"

	if err := smtp.SendMail(host+hostPort, auth, from, []string{to}, []byte(msg)); err != nil {
		return errors.Wrap(err, "error while sending message to email")
	}

	return nil
}

func sendWithPlayMobile(dmn *Daemon, text, phoneNumber, login, password, originator string) error {
	//[0] text [1] code
	textAndCode := strings.Split(text, ":")
	var code string
	if len(textAndCode) > 1 {
		code = textAndCode[1]
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		code = id.String()
	}
	cfg := config.Load()

	if originator == "" {
		originator = "3700"
	}

	var body models.Body
	client := http.Client{}

	message := models.Message{
		Recipient: phoneNumber,
		MessageID: fmt.Sprintf("%s%s", "del", code),
		SMS: models.SMS{
			Originator: originator,
			Content: models.Content{
				Text: text,
			},
		},
	}

	body.Messages = append(body.Messages, message)

	requestBody, err := json.Marshal(body)

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", dmn.Conf.PlayMobileUrl, bytes.NewBuffer(requestBody))

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	if login == "" || password == "" {
		login = cfg.PlayMobileLogin
		password = cfg.PlayMobilePassword
	}

	request.SetBasicAuth(login, password)

	res, err := client.Do(request)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("error while sending sms")
	}

	return nil
}

func (dmn *Daemon) sendWithMailchimp(text, to, from_email, key string) error {
	client := http.Client{}

	message := map[string]any{
		"key": key,
		"message": map[string]any{
			"from_email": from_email,
			"subject":    "OTP",
			"text":       text,
			"to": []map[string]any{
				{
					"email": to,
					"type":  "to",
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://mandrillapp.com/api/1.0/messages/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
