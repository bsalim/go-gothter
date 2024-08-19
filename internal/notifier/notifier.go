package notifier

import (
	"go-gothter/internal/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Notifier struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewNotifier(cfg *config.Config, logger *logrus.Logger) *Notifier {
	return &Notifier{cfg: cfg, logger: logger}
}

func (n *Notifier) SendNotification() {
	if n.cfg.Email.Enabled {
		// Create a new message
		m := gomail.NewMessage()
		m.SetHeader("From", n.cfg.Email.SMTPUser)
		m.SetHeader("To", n.cfg.Email.Recipient)
		m.SetHeader("Subject", n.cfg.Email.Subject)
		m.SetBody("text/plain", "Suspicious activity detected!")

		// Create a new dialer
		d := gomail.NewDialer(n.cfg.Email.SMTPServer, n.cfg.Email.SMTPPort, n.cfg.Email.SMTPUser, n.cfg.Email.SMTPPassword)

		// Send the email
		if err := d.DialAndSend(m); err != nil {
			n.logger.Errorf("Failed to send email: %v", err)
		} else {
			n.logger.Info("Email sent successfully!")
		}
	}
}
