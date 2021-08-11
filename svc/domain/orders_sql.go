package domain

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func insertFullOrder(ctx context.Context, db *pgxpool.Pool, order Order) (Order, error) {
	order, err := insertNewOrder(ctx, db, order)
	if err != nil {
		return order, err
	}

	// prepare the order id
	for i, pymt := range order.Payments {
		pymt.OrderID = order.ID

		var err error
		pymt, err = insertNewPayment(ctx, db, pymt)
		if err != nil {
			return order, err
		}
		order.Payments[i] = pymt
	}

	for i, li := range order.LineItems {
		li.OrderID = order.ID

		var err error
		li, err = insertNewLineItem(ctx, db, li)
		if err != nil {
			return order, err
		}
		order.LineItems[i] = li
	}

	order.Total.OrderID = order.ID
	total, err := insertNewOrderTotal(ctx, db, order.Total)
	if err != nil {
		return order, err
	}

	order.Total = total

	for i, m := range order.Metadata {
		m.OrderID = order.ID

		var err error
		m, err := insertOrderMetadata(ctx, db, m)
		if err != nil {
			return order, err
		}

		order.Metadata[i] = m
	}

	return order, nil
}

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
   INSERT INTO line_items (id, order_id, quantity, tax_amount, unit_price_amount, total_line_amount,
                           tax_amount_before_discount, unit_price_amount_before_discount,
                           total_line_amount_before_discount, unit_price_currency, tax_amount_currency, description,
                           created_at, updated_at)
   VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW()::timestamptz, NOW()::timestamptz)
RETURNING id, created_at, updated_at`

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

func insertNewOrderTotal(ctx context.Context, db *pgxpool.Pool, total Total) (Total, error) {
	const insertOrderTotalStmt = `
   INSERT INTO order_totals (id, order_id, amount, total_amount, tax_amount, total_amount_before_discount,
                             tax_amount_before_discount, created_at, updated_at)
   VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW()::timestamptz, NOW()::timestamptz)
RETURNING id, created_at, updated_at`
	row := db.QueryRow(ctx, insertOrderTotalStmt,
		total.OrderID,
		total.Amount,
		total.TotalAmount,
		total.TaxAmount,
		total.TotalAmountBeforeDiscount,
		total.TaxAmountBeforeDiscount,
	)

	if err := row.Scan(&total.ID, &total.CreatedAt, &total.UpdatedAt); err != nil {
		return total, DatabaseQueryError{
			Stmt:  insertOrderTotalStmt,
			Inner: err,
		}
	}

	return total, nil
}

func insertOrderMetadata(ctx context.Context, db *pgxpool.Pool, md Metadata) (Metadata, error) {
	const insertMetadataStmt = `
   INSERT INTO order_metadata (id, order_id, "key", "value", created_at, updated_at)
   VALUES (gen_random_uuid(), $1, $2, $3, now()::timestamptz, now()::timestamptz)
RETURNING id, created_at, updated_at`

	row := db.QueryRow(ctx, insertMetadataStmt, md.OrderID, md.Key, md.Value)

	if err := row.Scan(&md.ID, &md.CreatedAt, &md.UpdatedAt); err != nil {
		return md, DatabaseQueryError{
			Stmt:  insertMetadataStmt,
			Inner: err,
		}
	}

	return md, nil
}
