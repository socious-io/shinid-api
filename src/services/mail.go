package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"shin/src/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var SendGridTemplates map[string]string = map[string]string{
	"otp":             "d-0146441b623f4cb78833c50eb1a8c813",
	"forget-password": "d-d7aea3b78df042e8a2fdc83953960259",
}

var SendGridClient SendGridType

type SendGridType struct {
	ApiKey   string
	Url      string
	Disabled bool
}

func (sgc *SendGridType) SendWithTemplate(address string, name string, templateId string, items map[string]string) error {
	if sgc.Disabled {
		return nil
	}
	//Create Mail payload
	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail("Socious", "no-replay@socious.io"))
	m.SetTemplateID(templateId)

	//Adding Personalization
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(name, address),
	}
	p.AddTos(tos...)
	for key, value := range items {
		p.SetDynamicTemplateData(key, value)
	}
	m.AddPersonalizations(p)

	//Setup the request
	request := sendgrid.GetRequest(sgc.ApiKey, "/v3/mail/send", sgc.Url)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
		return err
	} else if strings.Split(strconv.Itoa(response.StatusCode), "")[0] != "2" {
		return errors.New(response.Body)
	}
	return nil
}

func InitSendGridService() {
	SendGridClient = SendGridType{
		Disabled: config.Config.Sendgrid.Disabled,
		ApiKey:   config.Config.Sendgrid.ApiKey,
		Url:      config.Config.Sendgrid.URL,
	}
}
