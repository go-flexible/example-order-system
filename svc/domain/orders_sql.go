package domain

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func insertNewOrder(ctx context.Context, db *pgxpool.Pool, order Order) (Order, error) {
	const insertStmt = `
   INSERT INTO orders (id, number, created_at, updated_at)
   VALUES (gen_random_uuid(), $1, now()::timestamptz, now()::timestamptz)
RETURNING id, number, created_at, updated_at;`

	row := db.QueryRow(ctx, insertStmt, order.OrderNumber)

	if err := row.Scan(&order.ID, &order.OrderNumber, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return order, DatabaseQueryError{
			Stmt:  insertStmt,
			Inner: err,
		}
	}
	return order, nil
}

func insertNewPayment(ctx context.Context, db *pgxpool.Pool, payment Payment) (Payment, error) {
	const insertPaymentStmt = `
   INSERT INTO payments (id, order_id, payment_method, amount, tax_amount, total_amount, currency, created_at, updated_at)
   VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, now()::timestamptz ,now()::timestamptz)
RETURNING id, created_at, updated_at`

	row := db.QueryRow(ctx, insertPaymentStmt,
		payment.OrderID, payment.PaymentMethod, payment.Amount,
		payment.TaxAmount, payment.TotalAmount, payment.Currency,
	)

	if err := row.Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
		return Payment{}, DatabaseQueryError{
			Stmt:  insertPaymentStmt,
			Inner: err,
		}
	}
	return payment, nil
}

func insertNewLineItem(ctx context.Context, db *pgxpool.Pool, li LineItem) (LineItem, error) {
	const insertLineItemStmt = `
	INSERT INTO line_items (
		id,
		order_id,
		quantity,
		tax_amount,
		unit_price_amount,
		total_line_amount,
		tax_amount_before_discount,
		unit_price_amount_before_discount,
		total_line_amount_before_discount,
		unit_price_currency,
		tax_amount_currency,
		description,
		created_at,
		updated_at
	)
	VALUES (
		gen_random_uuid(),
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		now()::timestamptz,
		now()::timestamptz
	)
	RETURNING id, created_at, updated_at
	`

	row := db.QueryRow(ctx, insertLineItemStmt,
		li.OrderID,
		li.Quantity,
		li.TaxAmount,
		li.UnitPriceAmount,
		li.TotalLineAmount,
		li.TaxAmountBeforeDiscount,
		li.UnitPriceAmountBeforeDiscount,
		li.TotalLineAmountBeforeDiscount,
		li.UnitPriceCurrency,
		li.TaxAmountCurrency,
		li.Description,
	)

	if err := row.Scan(&li.ID, &li.CreatedAt, &li.UpdatedAt); err != nil {
		return li, DatabaseQueryError{
			Stmt:  insertLineItemStmt,
			Inner: err,
		}
	}

	return li, nil
}
