package gmail

import (
	"fmt"
	"html/template"
	"strings"

	"main/internal/domain/contracts"
	"main/internal/infrastructure/requester/models"
	senderModels "main/internal/infrastructure/sender/models"
)

type requester struct {
	sender         contracts.IEmailSender
	ownerEmail     string
	reqestTmplPath string
	signature      string
}

func NewRequester(
	sender contracts.IEmailSender,
	ownerEmail,
	signature,
	baseTmplPath string,
) *requester {
	return &requester{
		sender:         sender,
		ownerEmail:     ownerEmail,
		signature:      signature,
		reqestTmplPath: baseTmplPath + "doctor_reminder.tmpl",
	}
}

func (n *requester) Request(subject, content, recipientEmail string) error {
	tmplData := models.BaseRequest{
		Subject:   subject,
		Content:   content,
		Signature: n.signature,
	}

	tmpl, err := template.ParseFiles(n.reqestTmplPath)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	err = n.send(recipientEmail, subject, tmpl, tmplData)
	if err != nil {
		return fmt.Errorf("could not send msg: %w", err)
	}

	return nil
}

func (n *requester) send(
	recipientEmail string,
	subject string,
	tmpl *template.Template,
	tmplData any,
) error {
	msg, err := n.buildMsg(recipientEmail, subject, tmpl, tmplData)
	if err != nil {
		return fmt.Errorf("could not build msg to recipient %s: %w", recipientEmail, err)
	}

	if err = n.sender.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg: %w", err)
	}

	return nil
}

func (n *requester) buildMsg(
	recipientEmail,
	subject string,
	tmpl *template.Template,
	tmplData any,
) (senderModels.Message, error) {
	var body strings.Builder

	err := tmpl.Execute(&body, tmplData)
	if err != nil {
		return senderModels.Message{}, fmt.Errorf("could not execute template: %w", err)
	}

	return senderModels.Message{
		From:    n.ownerEmail,
		To:      recipientEmail,
		Subject: subject,
		Body:    body.String(),
	}, nil
}
