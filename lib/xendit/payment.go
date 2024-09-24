package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
	xendit "github.com/xendit/xendit-go/v5"

	"github.com/xendit/xendit-go/v5/common"
	"github.com/xendit/xendit-go/v5/invoice"
)

type XenditClient struct {
	Key string
}

func (x XenditClient) CreateInvoice(amount int64, exCode, payerEmail, desc string, userFullname string, userPhoneNumber string) (*invoice.Invoice, *common.XenditSdkError) {
	client := xendit.NewClient(x.Key)
	paymentChannel := viper.GetStringSlice("xendit.available_channel")
	successUrl := viper.GetString("xendit.success_url")
	failureUrl := viper.GetString("xendit.failure_url")
	invoiceDuration := "600" // 10 minutes

	var (
		sendEmail = true
	)

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(exCode, float64(amount))
	createInvoiceRequest.PayerEmail = &payerEmail
	createInvoiceRequest.Description = &desc
	createInvoiceRequest.PaymentMethods = paymentChannel
	createInvoiceRequest.ShouldSendEmail = &sendEmail
	createInvoiceRequest.SuccessRedirectUrl = &successUrl
	createInvoiceRequest.FailureRedirectUrl = &failureUrl
	createInvoiceRequest.InvoiceDuration = &invoiceDuration

	customer := *invoice.NewCustomerObject()
	customer.GivenNames.Set(&userFullname)
	customer.PhoneNumber.Set(&userPhoneNumber)
	customer.MobileNumber.Set(&userPhoneNumber)
	customer.Email.Set(&payerEmail)
	createInvoiceRequest.SetCustomer(customer)

	resp, httpResponse, err := client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", httpResponse)
	}

	return resp, err
}

func IsCallbackTokenVerified(token string) bool {
	callbackToken := viper.GetString("xendit.callback_token")

	return callbackToken == token
}
