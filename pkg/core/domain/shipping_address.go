package domain

type ShippingAddress struct {
	UserId int64 `json:"user_id"`

	EmailInvoice string `json:"email_invoice"`
	FullName     string `json:"full_name"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state,omitempty"`
	PostCode     string `json:"post_code"`
	Country      string `json:"country"`
}
