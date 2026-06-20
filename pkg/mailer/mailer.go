package mailer

type Mailer interface {
	Send(to, subject, body string) error
}

type DevMailer struct{}

func (m *DevMailer) Send(to, subject, body string) error {
	return nil
}
