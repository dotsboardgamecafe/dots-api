package point_of_sale

import (
	SubModule "dots-api/lib/point_of_sale/sub_modules"
	"errors"
	"time"

	"github.com/spf13/viper"
)

type IPointOfSale interface {
	AllowToRedeemTheSameInvoice() bool
	GetInvoices(invoiceCode string) (map[string]interface{}, error)
	GetInvoice(invoiceCode string) (map[string]interface{}, SubModule.CloseOrderDetail, error)
	GetProductDetail(productId int64, accessToken string) (map[string]interface{}, SubModule.ProductDetail, error)
	GetInvoiceCodeFromList(invoiceCode string) (map[string]interface{}, SubModule.CloseOrderDetail, error)
}

func GetPointOfSale(pos string) (IPointOfSale, error) {
	if pos == "Olsera" {
		keys := []SubModule.OlseraKey{}

		for k, v := range viper.Get("olsera_pos").(map[string]interface{}) {
			var expiresAt time.Time
			if expiresAtString, ok := v.(map[string]interface{})["expires_at"].(string); ok {
				expiresAt, _ = time.Parse(time.RFC3339, expiresAtString)
			} else {
				expiresAt = v.(map[string]interface{})["expires_at"].(time.Time)
			}

			keys = append(keys, SubModule.OlseraKey{
				Name:             k,
				EnableRedeemOnce: int(v.(map[string]interface{})["enable_redeem_once"].(float64)),
				AppId:            v.(map[string]interface{})["app_id"].(string),
				SecretKey:        v.(map[string]interface{})["secret_key"].(string),
				AccessToken:      v.(map[string]interface{})["access_token"].(string),
				ExpiresAt:        expiresAt,
			})
		}

		return &SubModule.Olsera{
			Keys: keys,
		}, nil
	}

	return nil, errors.New("invalid POS type passed")
}
