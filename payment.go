package gopay

import (
	"context"
	"fmt"

	paymentApi "github.com/tkliner/go-gopay/apis/payment"
	"github.com/tkliner/go-gopay/client"
)

type PaymentGetter interface {
	Payment() PaymentInterface
}

type PaymentInterface interface {
	GetPayment(ctx context.Context, id int64) (payment *paymentApi.PaymentResponse, err error)
}

type payment struct {
	client client.Interface
}

func newPayment(c client.Interface) PaymentInterface {
	return &payment{
		client: c,
	}
}

func (p *payment) GetPayment(ctx context.Context, id int64) (payment *paymentApi.PaymentResponse, err error) {
	resp := &paymentApi.PaymentResponse{}	
	req := p.client.Get().Resource(fmt.Sprintf("%s/%d", pathPayment, id))
	
	err = req.Do(ctx).Convert(resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}