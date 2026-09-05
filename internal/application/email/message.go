package email

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

type DoctorReminder struct {
	Subject         string
	Content         string
	RecipientEmail  string
	SenderEmail     string
	SenderSignature string
}

type DoctorReminderTmplData struct {
	Content   string
	Signature string
}

type CeneoCatcher struct {
	ProductName     string
	ProductPrice    string
	ProductCompany  string
	ProductURL      string
	RecipientEmail  string
	SenderEmail     string
	SenderSignature string
}

type CeneoCatcherTmplData struct {
	Content   string
	Signature string
}
