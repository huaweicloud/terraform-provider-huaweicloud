package cbc

import (
	"strings"

	"github.com/chnsz/golangsdk"
)

// PaySubscriptionOrder is a method used to pay for a subscription order.
func PaySubscriptionOrder(client *golangsdk.ServiceClient, orderId string) error {
	httpUrl := "v3/orders/customer-orders/pay"
	payPath := client.Endpoint + httpUrl
	payPath = strings.ReplaceAll(payPath, "{project_id}", client.ProjectID)

	payOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildPaySubscriptionOrderBodyParams(orderId),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", payPath, &payOpts)
	return err
}

func buildPaySubscriptionOrderBodyParams(orderId string) map[string]interface{} {
	// `use_coupon`: Whether to use a coupon, the valid values are "YES" and "NO".
	// `use_discount`: Whether to use a discount, the valid values are "YES" and "NO".
	return map[string]interface{}{
		"order_id":     orderId,
		"use_coupon":   "NO",
		"use_discount": "NO",
	}
}
