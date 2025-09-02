package gopay

import (
	"context"
	"fmt"
	"testing"

	"github.com/tkliner/go-gopay/client/config"
)

func TestMock( _ *testing.T) {
	cfg := config.NewConfig(
		config.WithGatewayURL("https://gw.sandbox.gopay.com"),
		config.WithCredentials(8836046164, "1253288454", "Cdf5ChEA"),
	)

	client, err := New(cfg)

	if err != nil {
		panic(err)
	}

	resp, err := client.Payment().GetPayment(context.TODO(), 3283981064)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %+v\n", resp)
}