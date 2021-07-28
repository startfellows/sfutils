package sfutils

import (
	"context"
	"github.com/prometheus/common/log"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

var pathCredentials string
var packageName string
var err error
var service *androidpublisher.Service

func init() {
	ctx := context.Background()
	service, err = androidpublisher.NewService(ctx, option.WithCredentialsFile(pathCredentials))
	if err != nil {
		log.Fatalf("create service: %s", err)
		panic(err)
	}
}

func CheckGooglePlay(productId, token string) (*androidpublisher.ProductPurchase, error) {
	product, err := service.Purchases.Products.Get(packageName, productId, token).Do()
	return product, err
}
