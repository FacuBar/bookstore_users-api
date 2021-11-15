package domain

type PaymentOption struct {
	UserID int64 `json:"user_id"`

	Id          int64  `json:"id"`
	CardType    string `json:"card_type"`
	CardNumber  string `json:"card_number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	NameOnCard  string `json:"name_on_card"`
	CVV         string `json:"cvv"`

	BillingAddress ShippingAddress `json:"billing_address"`
}
