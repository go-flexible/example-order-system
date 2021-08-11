CREATE TABLE IF NOT EXISTS orders (
	id UUID NOT NULL,
	number STRING NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id ASC),
	UNIQUE INDEX orders_number_uindex (number ASC),
	FAMILY "primary" (id, number, created_at, updated_at)
)