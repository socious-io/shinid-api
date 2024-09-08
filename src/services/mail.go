package services

import (
	"fmt"
	"shin/src/lib"
	"shin/src/utils"
)

const EmailChannel = "email"

type EmailConfig struct {
	Approach    string // [template|direct]
	Destination string
	Title       string
	Template    string
	Args        map[string]string
}

func SendEmail(emailConfig EmailConfig) {
	Mq.sendJson(EmailChannel, emailConfig)
}

func EmailWorker(message interface{}) {
	emailConfig := new(EmailConfig)
	utils.Copy(message, emailConfig)

	var (
		destination = emailConfig.Destination
		title       = emailConfig.Title
		template    = lib.SendGridTemplates[emailConfig.Template]
		args        = emailConfig.Args
	)

	if emailConfig.Approach == "template" {
		//Sending email with template
		err := lib.SendGridClient.SendWithTemplate(destination, title, template, args)
		if err != nil {
			fmt.Println("Coudn't Send Email, Error: ", err.Error())
		}
	}
}
