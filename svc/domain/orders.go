package domain

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
)

type OrderController struct {
	Dependencies
	Nats NatsPublisher
}

func (o OrderController) New(ctx context.Context, order Order) error {
	order.OrderNumber = newOrderNumber(0)

	var err error
	order, err = insertFullOrder(ctx, o.DB(), order)
	if err != nil {
		return err
	}

	err = o.Nats.Publish("orders.new", order)
	if err != nil {
		return PublishingError{error: err}
	}

	return nil
}

const defaultOrderNumberLength = 12

func newOrderNumber(length int) string {
	const (
		alpha    = "BCDEFGHJKLMNPQRSTUVWXYZ"
		num      = "12345679"
		alphanum = alpha + num
	)

	if length == 0 {
		length = defaultOrderNumberLength
	}

	var id string

	for {
		if len(id) == length {
			return id
		}

		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphanum))))
		if err != nil {
			// at this point if we can't read random bytes,
			// you may as well give up!
			log.Fatal(err)
		}

		n := num.Int64()
		id += string(alphanum[n])
	}
}
