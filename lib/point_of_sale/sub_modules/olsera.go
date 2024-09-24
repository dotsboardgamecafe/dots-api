package sub_modules

import (
	"dots-api/lib/utils"
	"dots-api/services/api/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	BASE_URL = "https://api-open.olsera.co.id/api/open-api/"
	VERSION  = "v1"
	LANG     = "id"
)

type (
	Olsera struct {
		AppId     string
		SecretKey string
	}

	BearerToken struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	ErrorResponse struct {
		Error any `json:"error"`
	}

	MainError struct {
		Error SubError `json:"error"`
	}

	SubError struct {
		Id           string `json:"id"`
		Message      string `json:"message"`
		StatusCode   int    `json:"status_code"`
		ErrorMessage string `json:"error"`
	}

	CloseOrder struct {
		Data []struct {
			Id      int64  `json:"id"`
			OrderNo string `json:"order_no"`
		} `json:"data"`
	}

	OrderItems struct {
		Price             float64 `json:"price"`
		ProductId         int64   `json:"product_id"`
		Name              string  `json:"product_name"`
		SKU               string  `json:"product_sku"`
		CategoryId        int64   `json:"category_id"`
		CategoryName      string  `json:"category_name"`
		ClasificationId   int64   `json:"klasifikasi_id"`
		ClasificationName string  `json:"klasifikasi"`
		Quantity          int     `json:"qty"`
	}

	CloseOrderDetail struct {
		Data struct {
			Id           int64        `json:"id"`
			OrderNo      string       `json:"order_no"`
			Status       string       `json:"status_desc"`
			OrderAmount  string       `json:"order_amount"`
			TotalAmount  string       `json:"total_amount"`
			IsPaid       int          `json:"is_paid"`
			TotalItemQty int          `json:"total_item_qty"`
			Items        []OrderItems `json:"orderitems"`
			CreatedTime  string       `json:"created_time"`
		} `json:"data"`
	}

	ProductDetail struct {
		Data struct {
			Id                int64  `json:"id"`
			Name              string `json:"name"`
			SKU               string `json:"sku"`
			ClasificationId   int64  `json:"klasifikasi_id"`
			ClasificationName string `json:"klasifikasi"`
			CategoryId        int64  `json:"category_id"`
			CategoryName      string `json:"category_name"`
			Description       string `json:"description"`
		} `json:"data"`
	}
)

func (pos *Olsera) GenerateAccessToken() (string, error) {
	var (
		errorCallback ErrorResponse
		mainError     MainError

		callback BearerToken
	)
	url := BASE_URL + VERSION + "/" + LANG + "/token"

	values := map[string]interface{}{
		"app_id":     pos.AppId,
		"secret_key": pos.SecretKey,
		"grant_type": "secret_key",
	}

	request, err := utils.RequestHandler(values, url, http.MethodPost)
	if err != nil {
		return callback.AccessToken, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := utils.ResponseAsyncHandler(request)

	json.Unmarshal([]byte(response[0]), &errorCallback)
	if errorCallback.Error != nil || err != nil {
		json.Unmarshal([]byte(response[0]), &mainError)
		return callback.AccessToken, errors.New(mainError.Error.Message)
	}

	json.Unmarshal([]byte(response[0]), &callback)
	utils.ResponseHandler(request)

	return callback.AccessToken, err
}

func (pos *Olsera) AllowToRedeemTheSameInvoice() bool {
	return viper.GetInt("olsera_pos.enable_redeem_once") == 0
}

func (pos *Olsera) GetInvoices(invoiceCode string) (map[string]interface{}, error) {
	var (
		errorCallback ErrorResponse
		mainError     MainError

		callback  map[string]interface{}
		emptyBody map[string]interface{}
	)
	url := BASE_URL + VERSION + "/" + LANG + "/order/closeorder"

	if invoiceCode != "" {
		url += "?search=" + invoiceCode
	}

	accessToken, _ := pos.GenerateAccessToken()

	request, err := utils.RequestHandler(emptyBody, url, http.MethodGet)
	if err != nil {
		return callback, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := utils.ResponseAsyncHandler(request)

	json.Unmarshal([]byte(response[0]), &errorCallback)
	if errorCallback.Error != float64(0) || err != nil {
		json.Unmarshal([]byte(response[0]), &mainError)
		return callback, errors.New(mainError.Error.Message)
	}

	json.Unmarshal([]byte(response[0]), &callback)
	utils.ResponseHandler(request)

	return callback, err
}

func (pos *Olsera) GetInvoice(invoiceCode string) (map[string]interface{}, *CloseOrderDetail, error) {
	var (
		errorCallback ErrorResponse
		mainError     MainError

		orderDetailTrx CloseOrderDetail

		callback  map[string]interface{}
		emptyBody map[string]interface{}
	)

	url := BASE_URL + VERSION + "/" + LANG + "/order/closeorder/detail"

	if invoiceCode != "" {
		url += "?id=" + invoiceCode
	}

	accessToken, _ := pos.GenerateAccessToken()

	request, err := utils.RequestHandler(emptyBody, url, http.MethodGet)
	if err != nil {
		return callback, nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := utils.ResponseAsyncHandler(request)

	json.Unmarshal([]byte(response[0]), &errorCallback)
	if errorCallback.Error != float64(0) || err != nil {
		json.Unmarshal([]byte(response[0]), &mainError)
		return emptyBody, nil, errors.New(mainError.Error.Message)
	}

	json.Unmarshal([]byte(response[0]), &orderDetailTrx)

	productItems := orderDetailTrx.Data.Items
	for i := 0; i < len(productItems); i++ {
		item := &productItems[i]

		_, selectedProduct, _ := pos.GetProductDetail(item.ProductId, accessToken)

		item.CategoryId = selectedProduct.Data.CategoryId
		item.CategoryName = selectedProduct.Data.CategoryName
	}

	utils.ResponseHandler(request)

	return nil, &orderDetailTrx, err
}

func (pos *Olsera) GetProductDetail(productId int64, accessToken string) (map[string]interface{}, *ProductDetail, error) {
	var (
		errorCallback ErrorResponse
		mainError     MainError

		productDetail ProductDetail

		callback  map[string]interface{}
		emptyBody map[string]interface{}
	)

	url := BASE_URL + VERSION + "/" + LANG + "/product/detail?id=" + fmt.Sprint(productId)

	if accessToken == "" {
		accessToken, _ = pos.GenerateAccessToken()
	}

	request, err := utils.RequestHandler(emptyBody, url, http.MethodGet)
	if err != nil {
		return callback, nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := utils.ResponseAsyncHandler(request)

	json.Unmarshal([]byte(response[0]), &errorCallback)
	if errorCallback.Error != float64(0) || err != nil {
		json.Unmarshal([]byte(response[0]), &mainError)
		return emptyBody, nil, errors.New(mainError.Error.Message)
	}

	json.Unmarshal([]byte(response[0]), &productDetail)
	utils.ResponseHandler(request)

	return nil, &productDetail, err
}

func (pos *Olsera) GetInvoiceCodeFromList(invoiceCode string) (*map[string]interface{}, *CloseOrderDetail, error) {
	var (
		callback  map[string]interface{}
		emptyBody map[string]interface{}

		errorCallback ErrorResponse
		mainError     MainError

		orderTrx       CloseOrder
		orderDetailTrx CloseOrderDetail
	)
	orderUri := BASE_URL + VERSION + "/" + LANG + "/order/closeorder?search=" + invoiceCode

	accessToken, _ := pos.GenerateAccessToken()

	request, err := utils.RequestHandler(emptyBody, orderUri, http.MethodGet)
	if err != nil {
		return &callback, nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := utils.ResponseAsyncHandler(request)

	json.Unmarshal([]byte(response[0]), &errorCallback)
	if errorCallback.Error != float64(0) || err != nil {
		json.Unmarshal([]byte(response[0]), &mainError)
		return &emptyBody, nil, errors.New(mainError.Error.Message)
	}

	json.Unmarshal([]byte(response[0]), &orderTrx)
	getFirstOrder := orderTrx.Data[0]

	if getFirstOrder.Id != 0 {
		orderDetailUri := BASE_URL + VERSION + "/" + LANG + "/order/closeorder/detail?id=" + fmt.Sprint(getFirstOrder.Id)

		request, err := utils.RequestHandler(emptyBody, orderDetailUri, http.MethodGet)
		if err != nil {
			return &callback, nil, err
		}

		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		request.Header.Set("Authorization", "Bearer "+accessToken)
		response, err := utils.ResponseAsyncHandler(request)

		json.Unmarshal([]byte(response[0]), &errorCallback)
		if errorCallback.Error != float64(0) || err != nil {
			json.Unmarshal([]byte(response[0]), &mainError)
			return &emptyBody, nil, errors.New(mainError.Error.Message)
		}

		json.Unmarshal([]byte(response[0]), &orderDetailTrx)
	}

	utils.ResponseHandler(request)

	productItems := orderDetailTrx.Data.Items
	for i := 0; i < len(productItems); i++ {
		item := &productItems[i]

		_, selectedProduct, _ := pos.GetProductDetail(item.ProductId, accessToken)

		item.CategoryId = selectedProduct.Data.CategoryId
		item.CategoryName = selectedProduct.Data.CategoryName
		item.ClasificationId = selectedProduct.Data.ClasificationId
		item.ClasificationName = selectedProduct.Data.ClasificationName
	}

	return nil, &orderDetailTrx, err
}

/*
invoiceDetail module function
*/
func (invoiceDetail *CloseOrderDetail) GetTotalAmount() (float64, error) {
	return strconv.ParseFloat(invoiceDetail.Data.TotalAmount, 64)
}

func (invoiceDetail *CloseOrderDetail) GetLineOfProducts() string {
	var productItems []string
	for _, p := range invoiceDetail.Data.Items {
		item := fmt.Sprintf("%d %s", p.Quantity, p.Name)
		productItems = append(productItems, item)
	}

	// Join the formatted information using a separator
	result := strings.Join(productItems, "; ")

	return result
}

func (invoiceDetail *CloseOrderDetail) SavedInformationToJSONString() ([]byte, error) {
	invoiceData := invoiceDetail.Data
	orderItems := make([]model.UserOrderItems, 0)

	userClaimedInvoice := new(model.UserClaimedInvoice)
	userClaimedInvoice.OrderCreatedTime = invoiceData.CreatedTime
	userClaimedInvoice.OrderStatus = invoiceData.Status
	userClaimedInvoice.OrderId = invoiceData.Id
	userClaimedInvoice.OrderNo = invoiceData.OrderNo
	userClaimedInvoice.OrderTotalAmount = invoiceData.OrderAmount
	userClaimedInvoice.OrderTotalQty = invoiceData.TotalItemQty

	productItems := invoiceData.Items
	for i := 0; i < len(productItems); i++ {
		item := productItems[i]

		orderItems = append(orderItems, model.UserOrderItems{
			ProductId:          item.ProductId,
			ProductName:        item.Name,
			CategoryName:       item.CategoryName,
			ProductSKU:         item.SKU,
			ProductPrice:       item.Price,
			Qty:                item.Quantity,
			ClassificationName: item.ClasificationName,
		})
	}

	userClaimedInvoice.OrderItems = orderItems

	return json.Marshal(userClaimedInvoice)
}
