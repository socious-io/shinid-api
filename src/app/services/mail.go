package services

import (
	"fmt"

	"shin/src/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var SendGridTemplates map[string]string = map[string]string{
	"otp": "d-0146441b623f4cb78833c50eb1a8c813",
}

var SendGridClient SendGridType

type SendGridType struct {
	ApiKey string
	Url    string
}

func (sgc *SendGridType) SendWithTemplate(address string, name string, templateId string, items map[string]string) error {

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

	_, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func InitSendGridService() {
	SendGridClient = SendGridType{
		ApiKey: config.Config.Sendgrid.ApiKey,
		Url:    config.Config.Sendgrid.URL,
	}
}
