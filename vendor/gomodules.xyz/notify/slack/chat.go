package slack

import (
	"context"
	"errors"
	"github.com/nlopes/slack"
	"gomodules.xyz/envconfig"
	"gomodules.xyz/notify"
)

const UID = "slack"

type Options struct {
	AuthToken string   `envconfig:"AUTH_TOKEN" required:"true"`
	Channel   []string `envconfig:"CHANNEL"`
}

type client struct {
	opt  Options
	body string
}

var _ notify.ByChat = &client{}

func New(opt Options) *client {
	return &client{opt: opt}
}

func Default() (*client, error) {
	var opt Options
	err := envconfig.Process(UID, &opt)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func Load(loader envconfig.LoaderFunc) (*client, error) {
	var opt Options
	err := envconfig.Load(UID, &opt, loader)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func (c client) UID() string {
	return UID
}

func (c client) WithBody(body string) notify.ByChat {
	c.body = body
	return &c
}

func (c client) To(to string, cc ...string) notify.ByChat {
	c.opt.Channel = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.Channel) == 0 {
		return errors.New("missing to")
	}

	s := slack.New(c.opt.AuthToken)
	for _, channel := range c.opt.Channel {
		if _, _, err := s.PostMessageContext(
			context.TODO(),
			channel,
			slack.MsgOptionText(c.body, false)); err != nil {
			return err
		}
	}
	return nil
}
