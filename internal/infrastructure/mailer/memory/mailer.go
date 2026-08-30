package memory

import "main/internal/application/email"

type Mailer struct {
	storage *Storage
}

func NewMailer(
	storage *Storage,
) *Mailer {
	return &Mailer{
		storage: storage,
	}
}

func (s *Mailer) Send(messages ...email.Message) error {
	for _, m := range messages {
		s.storage.AddView("from: " + m.From + `<hr style="margin:5px 0;">`)
		s.storage.AddView("to: " + m.To + `<hr style="margin:5px 0;">`)
		s.storage.AddView("subject: " + m.Subject + `<hr style="margin:5px 0 40px 0;">`)
		s.storage.AddView(m.Body + `<hr style="margin:40px 0;">`)
	}

	return nil
}
