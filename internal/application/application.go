package application

import "main/internal/application/email"

type IMailer interface {
	Send(messages ...email.Message) error
}
