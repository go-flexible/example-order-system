package domain_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-flexible/example-order-system/svc/domain"
)

// This isn't really a test.
// It's just used to create an order schema we can use in HTTP requests.
func TestGeneratedOrder(t *testing.T) {
	order := domain.Order{
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
		LineItems: []domain.LineItem{
			{
				UpdatedAt:                     gofakeit.Date(),
				CreatedAt:                     gofakeit.Date(),
				TaxAmountCurrency:             gofakeit.CurrencyShort(),
				UnitPriceCurrency:             gofakeit.CurrencyShort(),
				TaxAmountBeforeDiscount:       gofakeit.Number(10, 1000),
				TotalLineAmountBeforeDiscount: gofakeit.Number(10, 1000),
				TotalLineAmount:               gofakeit.Number(10, 1000),
				UnitPriceAmount:               gofakeit.Number(10, 1000),
				TaxAmount:                     gofakeit.Number(10, 1000),
				Quantity:                      gofakeit.Number(10, 1000),
				UnitPriceAmountBeforeDiscount: gofakeit.Number(10, 1000),
			},
		},
		Payments: []domain.Payment{
			{
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
				Currency:      gofakeit.CurrencyShort(),
				PaymentMethod: "card",
				Amount:        gofakeit.Number(10, 1000),
				TaxAmount:     gofakeit.Number(10, 1000),
				TotalAmount:   gofakeit.Number(10, 1000),
			},
		},
		Metadata: []domain.Metadata{
			{
				Key:   "app_version",
				Value: gofakeit.AppVersion(),
			},
		},
		Total: domain.Total{
			CreatedAt:                 gofakeit.Date(),
			UpdatedAt:                 gofakeit.Date(),
			TaxAmount:                 gofakeit.Number(10, 1000),
			TotalAmount:               gofakeit.Number(10, 1000),
			TaxAmountBeforeDiscount:   gofakeit.Number(10, 1000),
			TotalAmountBeforeDiscount: gofakeit.Number(10, 1000),
			Amount:                    gofakeit.Number(10, 1000),
		},
	}

	raw, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(raw))
}

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
