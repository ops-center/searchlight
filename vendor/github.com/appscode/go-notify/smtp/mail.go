package smtp

import (
	"crypto/tls"

	"github.com/appscode/go-notify"
	"github.com/kelseyhightower/envconfig"
	gomail "gopkg.in/gomail.v2"
)

const Uid = "smtp"

type Options struct {
	Host               string   `envconfig:"HOST" required:"true"`
	Port               int      `envconfig:"PORT" required:"true"`
	InsecureSkipVerify bool     `envconfig:"INSECURE_SKIP_VERIFY"`
	Username           string   `envconfig:"USERNAME" required:"true"`
	Password           string   `envconfig:"PASSWORD" required:"true"`
	From               string   `envconfig:"FROM" required:"true"`
	To                 []string `envconfig:"TO" required:"true"`
}

type client struct {
	opt  Options
	mail *gomail.Message
	body string
	html bool
}

var _ notify.ByEmail = &client{}

func New(opt Options) *client {
	mail := gomail.NewMessage()
	mail.SetHeader("From", opt.From)
	mail.SetHeader("To", opt.To...)
	return &client{
		opt:  opt,
		mail: mail,
	}
}

func Default() (*client, error) {
	var opt Options
	err := envconfig.Process(Uid, &opt)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func (c *client) From(from string) notify.ByEmail {
	c.mail.SetHeader("From", from)
	return c
}

func (c *client) WithSubject(subject string) notify.ByEmail {
	c.mail.SetHeader("Subject", subject)
	return c
}
func (c *client) WithBody(body string) notify.ByEmail {
	c.body = body
	return c
}

func (c *client) WithTag(tag string) notify.ByEmail {
	return c
}

func (c *client) To(to string, cc ...string) notify.ByEmail {
	tos := append([]string{to}, cc...)
	c.mail.SetHeader("To", tos...)
	return c
}

func (c *client) Send() error {
	if c.html {
		c.mail.SetBody("text/html", c.body)
	} else {
		c.mail.SetBody("text/plain", c.body)
	}

	var d *gomail.Dialer
	if c.opt.Username != "" && c.opt.Password != "" {
		d = gomail.NewDialer(c.opt.Host, c.opt.Port, c.opt.Username, c.opt.Password)
	} else {
		d = &gomail.Dialer{Host: c.opt.Host, Port: c.opt.Port}
	}
	if c.opt.InsecureSkipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return d.DialAndSend(c.mail)
}

func (c *client) SendHtml() error {
	c.html = true
	return c.Send()
}
