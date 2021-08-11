package domain_test

import (
	"context"
	"testing"

	"github.com/go-flexible/example-order-system/svc/domain"
)

func TestOrderController_New(t *testing.T) {
	t.Cleanup(func() {
		dropAllTables(t)
	})

	ctrl := domain.OrderController{Dependencies: dependencies{}, Nats: &natsPulisherMock{}}

	err := ctrl.New(context.Background(), domain.Order{})
	if err != nil {
		t.Fatal(err)
	}
}
