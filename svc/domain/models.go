package domain

import "time"

// Order is a struct that contains the order information.
type Order struct {
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	ID          string     `json:"id" db:"id"`
	OrderNumber string     `json:"order_number" db:"order_number"`
	LineItems   []LineItem `json:"line_items" db:"line_items"`
	Payments    []Payment  `json:"payments" db:"payments"`
	Metadata    []Metadata `json:"metadata" db:"metadata"`
	Total       Total      `json:"total" db:"total"`
}

// LineItem is a line item in an order.
type LineItem struct {
	UpdatedAt                     time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
	OrderID                       string    `json:"order_id" db:"order_id"`
	Description                   string    `json:"description" db:"description"`
	TaxAmountCurrency             string    `json:"tax_amount_currency" db:"tax_amount_currency"`
	UnitPriceCurrency             string    `json:"unit_price_currency" db:"unit_price_currency"`
	ID                            string    `json:"id" db:"id"`
	TaxAmountBeforeDiscount       int       `json:"tax_amount_before_discount" db:"tax_amount_before_discount"`
	TotalLineAmountBeforeDiscount int       `json:"total_line_amount_before_discount" db:"total_line_amount_before_discount"`
	TotalLineAmount               int       `json:"total_line_amount" db:"total_line_amount"`
	UnitPriceAmount               int       `json:"unit_price_amount" db:"unit_price_amount"`
	TaxAmount                     int       `json:"tax_amount" db:"tax_amount"`
	Quantity                      int       `json:"quantity" db:"quantity"`
	UnitPriceAmountBeforeDiscount int       `json:"unit_price_amount_before_discount" db:"unit_price_amount_before_discount"`
}

// Payment is a payment method for an order.
type Payment struct {
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	ID            string    `json:"id" db:"id"`
	Currency      string    `json:"currency" db:"currency"`
	OrderID       string    `json:"order_id" db:"order_id"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	Amount        int       `json:"amount" db:"amount"`
	TaxAmount     int       `json:"tax_amount" db:"tax_amount"`
	TotalAmount   int       `json:"total_amount" db:"total_amount"`
}

// Total is a struct that contains the order total information.
type Total struct {
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
	ID                        string    `json:"id" db:"id"`
	OrderID                   string    `json:"order_id" db:"order_id"`
	TaxAmount                 int       `json:"tax_amount" db:"tax_amount"`
	TotalAmount               int       `json:"total_amount" db:"total_amount"`
	TaxAmountBeforeDiscount   int       `json:"tax_amount_before_discount" db:"tax_amount_before_discount"`
	TotalAmountBeforeDiscount int       `json:"total_amount_before_discount" db:"total_amount_before_discount"`
	Amount                    int       `json:"amount" db:"amount"`
}

type Metadata struct {
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	ID        string    `json:"-" db:"id"`
	OrderID   string    `json:"-" db:"order_id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
}
