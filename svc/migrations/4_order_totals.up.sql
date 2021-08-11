CREATE TABLE IF NOT EXISTS order_totals (
	id UUID NOT NULL,
	order_id UUID NOT NULL,
	amount INT8 NOT NULL,
	total_amount INT8 NOT NULL,
	tax_amount INT8 NOT NULL,
	total_amount_before_discount INT8 NOT NULL,
	tax_amount_before_discount INT8 NULL,
	CONSTRAINT order_totals_orders_id_fk FOREIGN KEY (order_id) REFERENCES orders(id),
	INDEX order_totals_auto_index_order_totals_orders_id_fk (order_id ASC),
	UNIQUE INDEX order_totals_id_uindex (id ASC),
	FAMILY "primary" (id, order_id, amount, total_amount, tax_amount, total_amount_before_discount, tax_amount_before_discount, rowid)
)