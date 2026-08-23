package email

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

type DoctorReminder struct {
	Subject        string
	Content        string
	RecipientEmail string
}

type DoctorReminderTmplData struct {
	Content   string
	Signature string
}

type CeneoCatcher struct {
	ProductName    string
	ProductPrice   string
	ProductCompany string
	ProductURL     string
}

type CeneoCatcherTmplData struct {
	Content string
}
