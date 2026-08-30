package email

import (
	"fmt"
	"html/template"
	"strings"
)

type IMailer interface {
	Send(messages ...Message) error
}

type Composer struct {
	mailer         IMailer
	ownerEmail     string
	ownerSignature string
	baseTmplPath   string
}

func NewComposer(
	mailer IMailer,
	ownerEmail,
	ownerSignature,
	baseTmplPath string,
) *Composer {
	return &Composer{
		mailer:         mailer,
		ownerEmail:     ownerEmail,
		ownerSignature: ownerSignature,
		baseTmplPath:   baseTmplPath,
	}
}

func (n *Composer) ComposeForDoctorReminder(data DoctorReminder) (Message, error) {
	tmplData := DoctorReminderTmplData{
		Content:   data.Content,
		Signature: n.ownerSignature,
	}

	tmplPath := buildTmplPath(n.baseTmplPath, "doctor_reminder")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse template: %w", err)
	}

	msg, err := n.buildDoctorReminderMsg(tmpl, tmplData, data.Subject, data.RecipientEmail)
	if err != nil {
		return Message{},
			fmt.Errorf("could not build msg for doctor_reminder with data %v: %w", tmplData, err)
	}

	return msg, nil
}

func (n *Composer) ComposeForCeneoCatcher(data CeneoCatcher) (Message, error) {
	tmplData := CeneoCatcherTmplData{
		Content: data.ProductURL,
	}

	subject := fmt.Sprintf("%s: %s (%s)", data.ProductName, data.ProductPrice, data.ProductCompany)

	tmplPath := buildTmplPath(n.baseTmplPath, "ceneo_catcher")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse template: %w", err)
	}

	msg, err := n.buildCeneoCatcherMsg(tmpl, tmplData, subject)
	if err != nil {
		return Message{},
			fmt.Errorf("could not build msg for ceneo catcher with data %v: %w", tmplData, err)
	}

	return msg, nil
}

func buildTmplPath(tmplPath, taskName string) string {
	return tmplPath + taskName + ".tmpl"
}

func (n *Composer) buildDoctorReminderMsg(
	tmpl *template.Template,
	data DoctorReminderTmplData,
	subject string,
	recipientEmail string,
) (Message, error) {
	var body strings.Builder

	err := tmpl.Execute(&body, data)
	if err != nil {
		return Message{}, fmt.Errorf("could not execute template: %w", err)
	}

	return Message{
		From:    n.ownerEmail,
		To:      recipientEmail,
		Subject: subject,
		Body:    body.String(),
	}, nil
}

func (n *Composer) buildCeneoCatcherMsg(
	tmpl *template.Template,
	data CeneoCatcherTmplData,
	subject string,
) (Message, error) {
	var body strings.Builder

	err := tmpl.Execute(&body, data)
	if err != nil {
		return Message{}, fmt.Errorf("could not execute template: %w", err)
	}

	return Message{
		From:    n.ownerEmail,
		To:      n.ownerEmail,
		Subject: subject,
		Body:    body.String(),
	}, nil
}
