CREATE TABLE IF NOT EXISTS payments (
	id UUID NOT NULL,
	order_id UUID NOT NULL,
	payment_method STRING NOT NULL,
	amount INT8 NOT NULL,
	tax_amount INT8 NOT NULL,
	total_amount INT8 NOT NULL,
	currency STRING NOT NULL,
	created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT payments_orders_id_fk FOREIGN KEY (order_id) REFERENCES orders(id),
	INDEX payments_auto_index_payments_orders_id_fk (order_id ASC),
	UNIQUE INDEX payments_id_uindex (id ASC),
	FAMILY "primary" (id, order_id, payment_method, amount, tax_amount, total_amount, currency, rowid, created_at, updated_at)
)