package request

//	InvoiceCallbackJson := `{
//	  "id" : "593f4ed1c3d3bb7f39733d83",
//	  "external_id" : "testing-invoice",
//	  "user_id" : "5848fdf860053555135587e7",
//	  "payment_method" : "RETAIL_OUTLET",
//	  "status" : "PAID",
//	  "merchant_name" : "Xendit",
//	  "amount" : 2000000,
//	  "paid_amount" : 2000000,
//	  "paid_at" : "2020-01-14T02:32:50.912Z",
//	  "payer_email" : "test@xendit.co",
//	  "description" : "Invoice webhook test",
//	  "created" : "2020-01-13T02:32:49.827Z",
//	  "updated" : "2020-01-13T02:32:50.912Z",
//	  "currency" : "IDR",
//	  "payment_channel" : "ALFAMART",
//	  "payment_destination" : "TEST815"
//	}`
type InvoiceCallbackRequest struct {
	ID                 string `json:"id"`
	ExternalID         string `json:"external_id"`
	UserID             string `json:"user_id"`
	PaymentMethod      string `json:"payment_method"`
	Status             string `json:"status"`
	MerchantName       string `json:"merchant_name"`
	Amount             int    `json:"amount"`
	PaidAmount         int    `json:"paid_amount"`
	PaidAt             string `json:"paid_at"`
	PayerEmail         string `json:"payer_email"`
	Desc               string `json:"description"`
	CreatedAt          string `json:"created"`
	UpdatedAt          string `json:"updated"`
	Currency           string `json:"currency"`
	PaymentChannel     string `json:"payment_channel"`
	PaymentDestination string `json:"payment_destination"`
}
