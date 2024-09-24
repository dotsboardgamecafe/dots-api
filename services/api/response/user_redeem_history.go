package response

import "dots-api/services/api/model"

type UserReedemHistoryRes struct {
	UserCode           string `json:"user_code,omitempty"`
	CustomId           string `json:"custom_id,omitempty"`
	PointEarned        int    `json:"point_earned,omitempty"`
	InvoiceCode        string `json:"invoice_code,omitempty"`
	InvoiceAmount      int    `json:"invoice_amount,omitempty"`
	InvoiceDescription string `json:"invoice_description,omitempty"`
	IsInvoiceNew       *bool  `json:"is_invoice_new,omitempty"`
	ClaimedDate        string `json:"claimed_date,omitempty"`
}

type UserClaimHistoryRes struct {
	InvoiceCode   string                 `json:"invoice_code,omitempty"`
	InvoiceAmount int                    `json:"invoice_amount,omitempty"`
	InvoiceItems  []model.UserOrderItems `json:"invoice_items"`
	ClaimedTime   string                 `json:"claimed_time,omitempty"`
	ClaimedDate   string                 `json:"claimed_date,omitempty"`
}
