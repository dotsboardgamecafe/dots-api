package point_of_sale

import (
	SubModule "dots-api/lib/point_of_sale/sub_modules"
	"fmt"

	"github.com/spf13/viper"
)

type IPointOfSale interface {
	AllowToRedeemTheSameInvoice() bool
	GetInvoices(invoiceCode string) (map[string]interface{}, error)
	GetInvoice(invoiceCode string) (map[string]interface{}, *SubModule.CloseOrderDetail, error)
	GetProductDetail(productId int64, accessToken string) (map[string]interface{}, *SubModule.ProductDetail, error)
	GetInvoiceCodeFromList(invoiceCode string) (*map[string]interface{}, *SubModule.CloseOrderDetail, error)
}

func GetPointOfSale(pos string) (IPointOfSale, error) {
	if pos == "Olsera" {
		return &SubModule.Olsera{
			AppId:     viper.GetString("olsera_pos.app_id"),
			SecretKey: viper.GetString("olsera_pos.secret_key"),
		}, nil
	}

	return nil, fmt.Errorf("invalid POS type passed")
}
