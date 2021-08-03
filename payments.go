package sfutils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/prometheus/common/log"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
)

var PathCredentials string
var PackageName string
var err error
var service *androidpublisher.Service

// IOSResponseData Apple's response structure
type IOSResponseData struct {
	Status  float64 `json:"status"`
	Receipt struct {
		InApp []struct {
			ProductID             string `json:"product_id"`
			OriginalTransactionID string `json:"original_transaction_id"`
		} `json:"in_app"`
	} `json:"receipt"`
}

func init() {
	ctx := context.Background()
	service, err = androidpublisher.NewService(ctx, option.WithCredentialsFile(PathCredentials)) // create server for send request
	if err != nil {
		log.Fatalf("create service: %s", err)
		panic(err)
	}
}

// CheckGooglePlay request for a product by product_id
func CheckGooglePlay(productId, token string) (*androidpublisher.ProductPurchase, error) {
	product, err := service.Purchases.Products.Get(PackageName, productId, token).Do()
	return product, err
}

// CheckApplePay request the AppStore for complete product information
func CheckApplePay(password, receipt string) (*IOSResponseData, error) {
	prodUrl := "https://buy.itunes.apple.com/verifyReceipt"
	testUrl := "https://sandbox.itunes.apple.com/verifyReceipt"

	var bytesBody *bytes.Buffer

	postBody, _ := json.Marshal(map[string]interface{}{
		"password":                 password,
		"receipt-data":             receipt,
		"exclude-old-transactions": true,
	})
	bytesBody = bytes.NewBuffer(postBody)
	result, err := RequestToPayPlatform(prodUrl, "ios", bytesBody)
	if err != nil {
		return nil, err
	}

	status := result.Status

	if status == float64(21007) { // 21007 this is test server
		bytesBody = bytes.NewBuffer(postBody)
		result, err = RequestToPayPlatform(testUrl, "ios", bytesBody)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// RequestToPayPlatform parsing the answer and issuing an answer on a given platform
func RequestToPayPlatform(url, platform string, data *bytes.Buffer) (*IOSResponseData, error) {
	resp, err := http.Post(url, "application/json", data)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if platform == "ios" {
		responseData := &IOSResponseData{}
		err = json.Unmarshal([]byte(body), &responseData)
		if err != nil {
			return nil, err
		}
		return responseData, nil
	}

	return nil, nil
}
