package domain

import "time"

// Order is a struct that contains the order information.
type Order struct {
	ID          string              `json:"id"`           // uuid
	OrderNumber string              `json:"order_number"` // human readable order number
	LineItems   []LineItem          `json:"line_items"`
	Payments    []Payment           `json:"payments"`
	Metadata    map[string][]string `json:"metadata"`
	Total       Total               `json:"total"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// LineItem is a line item in an order.
type LineItem struct {
	ID                            string    `json:"id"`       // uuid
	OrderID                       string    `json:"order_id"` // uuid
	Quantity                      int       `json:"quantity"`
	TaxAmount                     int       `json:"tax_amount"`
	UnitPriceAmount               int       `json:"unit_price_amount"`
	TotalLineAmount               int       `json:"total_line_amount"`
	TaxAmountBeforeDiscount       int       `json:"tax_amount_before_discount"`
	UnitPriceAmountBeforeDiscount int       `json:"unit_price_amount_before_discount"`
	TotalLineAmountBeforeDiscount int       `json:"total_line_amount_before_discount"`
	UnitPriceCurrency             string    `json:"unit_price_currency"`
	TaxAmountCurrency             string    `json:"tax_amount_currency"`
	Description                   string    `json:"description"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
}

// Payment is a payment method for an order.
type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	Amount        int       `json:"amount"`
	TaxAmount     int       `json:"tax_amount"`
	TotalAmount   int       `json:"total_amount"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Total is a struct that contains the order total information.
type Total struct {
	ID                        string `json:"id"`
	OrderID                   string `json:"order_id"`
	Amount                    int    `json:"amount"`
	TaxAmount                 int    `json:"tax_amount"`
	TotalAmount               int    `json:"total_amount"`
	TaxAmountBeforeDiscount   int    `json:"tax_amount_before_discount"`
	TotalAmountBeforeDiscount int    `json:"total_amount_before_discount"`
}
