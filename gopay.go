package gopay

import (
	"net/http"

	"github.com/tkliner/go-gopay/client"
	"github.com/tkliner/go-gopay/client/config"
	gopayHttp "github.com/tkliner/go-gopay/client/http"
	"github.com/tkliner/go-gopay/client/logger"
)

const (
	pathPayment = "/payments/payment"
)

type Clienter interface {
	Client() client.Interface
	PaymentGetter
}

type GoPay struct {
	client client.Interface
	logger logger.Logger
}

func New(config *config.Config) (Clienter, error) {
	copy := *config
	defaults(&copy)
	httpClient, err := gopayHttp.NewHTTPClient(&copy)
	if err != nil {
		return nil, err
	}

	return NewWithClient(&copy, httpClient)

}

func NewWithClient(config *config.Config, c *http.Client) (Clienter, error) {
	copy := *config
	defaults(&copy)

	cl, err := client.NewClient(&copy, c)

	if err != nil {
		return nil, err
	}

	return &GoPay{
		client: cl,
		logger: copy.Logger,
	}, nil

}

func (g *GoPay) Payment() PaymentInterface {
	return newPayment(g.client)
}

func defaults(cfg *config.Config) {
	if cfg.Logger == nil {
		cfg.Logger = logger.NewNoOpLogger()
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = config.DefaultTimeout
	}

	if cfg.Language == "" {
		cfg.Language = config.DefaultLanguage
	}

}

func (g *GoPay) Client() client.Interface {
	return g.client
}