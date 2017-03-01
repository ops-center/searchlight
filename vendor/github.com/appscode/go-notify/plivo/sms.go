package plivo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/go-notify"
	"github.com/kelseyhightower/envconfig"
)

const (
	Uid         = "plivo"
	urlTemplate = "https://api.plivo.com/v1/Account/%v/Message/"
)

type Options struct {
	AuthID    string   `envconfig:"AUTH_ID" required:"true"`
	AuthToken string   `envconfig:"AUTH_TOKEN" required:"true"`
	From      string   `envconfig:"FROM" required:"true"`
	To        []string `envconfig:"TO" required:"true"`
}

type messageSendParams struct {
	Src  string `json:"src,omitempty"`
	Dst  string `json:"dst,omitempty"`
	Text string `json:"text,omitempty"`
}

type messageSendResponseBody struct {
	Error string `json:"error"`
}

type client struct {
	opt    Options
	to     []string
	url    string
	params messageSendParams
}

var _ notify.BySMS = &client{}

func New(opt Options) *client {
	return &client{
		opt: opt,
		to:  opt.To,
		url: fmt.Sprintf(urlTemplate, opt.AuthID),
		params: messageSendParams{
			Src: opt.From,
		},
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

func (c *client) From(from string) notify.BySMS {
	c.params.Dst = from
	return c
}

func (c *client) WithBody(body string) notify.BySMS {
	c.params.Text = body
	return c
}

func (c *client) To(to string, cc ...string) notify.BySMS {
	c.to = append([]string{to}, cc...)
	return c
}

func (c *client) Send() error {
	httpClient := &http.Client{Timeout: time.Second * 10}

	params := c.params
	for _, dst := range c.to {
		params.Dst = dst
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(params); err != nil {
			return err
		}

		req, err := http.NewRequest("POST", c.url, buf)
		if err != nil {
			return err
		}

		req.SetBasicAuth(c.opt.AuthID, c.opt.AuthToken)
		req.Header.Add("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		respBody := &messageSendResponseBody{}
		if err = json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return err
		}
		if respBody.Error != "" {
			return errors.New(respBody.Error)
		}

		resp.Body.Close()
	}
	return nil
}
