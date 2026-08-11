package contracts

import "main/internal/infrastructure/sender/models"

type IEmailSender interface {
	Send(messages ...models.Message) error
}
