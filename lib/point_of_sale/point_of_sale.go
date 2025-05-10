package point_of_sale

import (
	SubModule "dots-api/lib/point_of_sale/sub_modules"
	"fmt"

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

		for _, v := range viper.Get("olsera_pos").([]interface{}) {
			keys = append(keys, SubModule.OlseraKey{
				Name:             v.(map[string]interface{})["name"].(string),
				EnableRedeemOnce: int(v.(map[string]interface{})["enable_redeem_once"].(float64)),
				AppId:            v.(map[string]interface{})["app_id"].(string),
				SecretKey:        v.(map[string]interface{})["secret_key"].(string),
			})
		}

		return &SubModule.Olsera{
			Keys: keys,
		}, nil
	}

	return nil, fmt.Errorf("invalid POS type passed")
}
