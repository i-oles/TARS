package email

import (
	"fmt"
	"html/template"
	"strings"
)

type Composer struct {
	baseTmplPath string
}

func NewComposer(
	baseTmplPath string,
) *Composer {
	return &Composer{
		baseTmplPath: baseTmplPath,
	}
}

func (n *Composer) ComposeForDoctorReminder(data DoctorReminder) (Message, error) {
	tmplPath := buildTmplPath(n.baseTmplPath, "doctor_reminder")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse template: %w", err)
	}

	tmplData := DoctorReminderTmplData{
		Content:   data.Content,
		Signature: data.SenderSignature,
	}

	var body strings.Builder

	err = tmpl.Execute(&body, tmplData)
	if err != nil {
		return Message{}, fmt.Errorf("could not execute template: %w", err)
	}

	msg := n.buildMsg(body.String(), data.Subject, data.RecipientEmail, data.SenderEmail)

	return msg, nil
}

func (n *Composer) ComposeForCeneoCatcher(data CeneoCatcher) (Message, error) {
	tmplPath := buildTmplPath(n.baseTmplPath, "ceneo_catcher")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse template: %w", err)
	}

	tmplData := CeneoCatcherTmplData{
		Content:   data.ProductURL,
		Signature: data.SenderSignature,
	}

	var body strings.Builder

	err = tmpl.Execute(&body, tmplData)
	if err != nil {
		return Message{}, fmt.Errorf("could not execute template: %w", err)
	}

	subject := fmt.Sprintf("%s: %s (%s)", data.ProductName, data.ProductPrice, data.ProductCompany)

	msg := n.buildMsg(body.String(), subject, data.RecipientEmail, data.SenderEmail)

	return msg, nil
}

func buildTmplPath(tmplPath, taskName string) string {
	return tmplPath + taskName + ".tmpl"
}

func (n *Composer) buildMsg(
	body, subject, recipientEmail, senderEmail string,
) Message {
	return Message{
		From:    senderEmail,
		To:      recipientEmail,
		Subject: subject,
		Body:    body,
	}
}
