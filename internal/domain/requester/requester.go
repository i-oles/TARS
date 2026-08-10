package requester

type IRequester interface {
	Request(subject, content, recipientEmail string) error
}
