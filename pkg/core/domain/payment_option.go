package domain

type PaymentOption struct {
	UserID int64

	Id          int64
	CardType    string
	CardNumber  string
	ExpiryMonth int
	ExpiryYear  int
	NameOnCard  string
	CVV         string

	BillingAddress ShippingAddress
}
