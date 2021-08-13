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

func getFullOrder(ctx context.Context, db *pgxpool.Pool, id string) (Order, error) {
	order, err := getOrderByID(ctx, db, id)
	if err != nil {
		return Order{}, nil
	}

	payments, err := getPaymentsForOrder(ctx, db, id)
	if err != nil {
		return Order{}, nil
	}
	order.Payments = payments

	lineItems, err := getLineItemsForOrder(ctx, db, id)
	if err != nil {
		return Order{}, nil
	}
	order.LineItems = lineItems

	metadata, err := getMetadataForOrder(ctx, db, id)
	if err != nil {
		return Order{}, nil
	}
	order.Metadata = metadata

	total, err := getTotalForOrder(ctx, db, id)
	if err != nil {
		return Order{}, nil
	}
	order.Total = total

	return order, nil
}

func getOrderByID(ctx context.Context, db *pgxpool.Pool, id string) (Order, error) {
	const query = `SELECT id, number, created_at, updated_at FROM orders WHERE id = $1`

	row := db.QueryRow(ctx, query, id)

	var order Order
	err := row.Scan(&order.ID, &order.OrderNumber, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return Order{}, DatabaseQueryError{
			Stmt:  query,
			Inner: err,
		}
	}

	return order, nil
}

func getPaymentsForOrder(ctx context.Context, db *pgxpool.Pool, id string) ([]Payment, error) {
	const query = `
SELECT id,
       order_id,
       payment_method,
       amount,
       tax_amount,
       total_amount,
       currency,
       created_at,
       updated_at
  FROM payments
 WHERE order_id = $1`

	rows, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, DatabaseQueryError{
			Stmt:  query,
			Inner: nil,
		}
	}

	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.PaymentMethod,
			&p.Amount,
			&p.TaxAmount,
			&p.TotalAmount,
			&p.Currency,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, DatabaseQueryError{
				Stmt:  query,
				Inner: err,
			}
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func getLineItemsForOrder(ctx context.Context, db *pgxpool.Pool, id string) ([]LineItem, error) {
	const query = `
SELECT id,
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
  FROM line_items
 WHERE order_id = $1`

	rows, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, DatabaseQueryError{
			Stmt:  query,
			Inner: nil,
		}
	}

	defer rows.Close()

	var lineItems []LineItem
	for rows.Next() {
		var li LineItem
		err := rows.Scan(
			&li.ID,
			&li.OrderID,
			&li.Quantity,
			&li.TaxAmount,
			&li.UnitPriceAmount,
			&li.TotalLineAmount,
			&li.TaxAmountBeforeDiscount,
			&li.UnitPriceAmountBeforeDiscount,
			&li.TotalLineAmountBeforeDiscount,
			&li.UnitPriceCurrency,
			&li.TaxAmountCurrency,
			&li.Description,
			&li.CreatedAt,
			&li.UpdatedAt,
		)
		if err != nil {
			return nil, DatabaseQueryError{
				Stmt:  query,
				Inner: err,
			}
		}
		lineItems = append(lineItems, li)
	}

	return lineItems, nil
}

func getMetadataForOrder(ctx context.Context, db *pgxpool.Pool, id string) ([]Metadata, error) {
	const query = `
SELECT id,
       order_id,
       key,
       value,
       created_at,
       updated_at
  FROM order_metadata
 WHERE order_id = $1`

	rows, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, DatabaseQueryError{
			Stmt:  query,
			Inner: nil,
		}
	}

	defer rows.Close()

	var metadata []Metadata
	for rows.Next() {
		var md Metadata
		err := rows.Scan(
			&md.ID,
			&md.OrderID,
			&md.Key,
			&md.Value,
			&md.CreatedAt,
			&md.UpdatedAt,
		)
		if err != nil {
			return nil, DatabaseQueryError{
				Stmt:  query,
				Inner: err,
			}
		}
		metadata = append(metadata, md)
	}

	return metadata, nil
}

func getTotalForOrder(ctx context.Context, db *pgxpool.Pool, id string) (Total, error) {
	const query = `
SELECT id,
       order_id,
       amount,
       tax_amount,
       total_amount,
       tax_amount_before_discount,
       total_amount_before_discount,
       created_at,
       updated_at
  FROM order_totals
 WHERE order_id = $1`

	row := db.QueryRow(ctx, query, id)

	var total Total
	err := row.Scan(
		&total.ID,
		&total.OrderID,
		&total.Amount,
		&total.TaxAmount,
		&total.TotalAmount,
		&total.TaxAmountBeforeDiscount,
		&total.TotalAmountBeforeDiscount,
		&total.CreatedAt,
		&total.UpdatedAt,
	)
	if err != nil {
		return Total{}, DatabaseQueryError{
			Stmt:  query,
			Inner: err,
		}
	}

	return total, nil
}
