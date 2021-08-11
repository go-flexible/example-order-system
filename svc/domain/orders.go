package domain

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
)

type OrdersController struct{ Dependencies }

func (o OrdersController) New(ctx context.Context, order Order) error {
	order.OrderNumber = newOrderNumber(0)
	order, err := insertNewOrder(ctx, o.DB(), order)
	if err != nil {
		return err
	}

	// prepare the order id
	for i, pymt := range order.Payments {
		pymt.OrderID = order.ID

		var err error
		pymt, err = insertNewPayment(ctx, o.DB(), pymt)
		if err != nil {
			return err
		}
		order.Payments[i] = pymt
	}

	for i, li := range order.LineItems {
		li.OrderID = order.ID

		var err error
		li, err = insertNewLineItem(ctx, o.DB(), li)
		if err != nil {
			return err
		}
		order.LineItems[i] = li
	}

	order.Total.OrderID = order.ID
	total, err := insertNewOrderTotal(ctx, o.DB(), order.Total)
	if err != nil {
		return err
	}

	order.Total = total

	err = o.Nats().Publish("orders.new", order)
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
