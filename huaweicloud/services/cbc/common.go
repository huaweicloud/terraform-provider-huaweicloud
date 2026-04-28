package cbc

import (
	"fmt"
	"strings"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	// "github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/bss/v2/resources"
	// "github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	// "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	// "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	// "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
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

func UpdateAutoRenew(c *golangsdk.ServiceClient, enabled, resourceId string) error {
	if enabled == "true" {
		return resources.EnableAutoRenew(c, resourceId)
	}
	return resources.DisableAutoRenew(c, resourceId)
}

// GetResourceIDsByOrder returns resource IDs from an order.
func GetResourceIDsByOrder(client *golangsdk.ServiceClient, orderId string, onlyMainResource int) ([]string, error) {
	if strings.TrimSpace(orderId) == "" {
		return nil, fmt.Errorf("order id is empty")
	}
	listOpts := resources.ListOpts{
		OrderId:          orderId,
		OnlyMainResource: onlyMainResource,
	}
	resp, err := resources.List(client, listOpts)
	if err != nil {
		return nil, fmt.Errorf("error getting order (%s) details: %s", orderId, err)
	}
	if resp == nil || resp.TotalCount < 1 {
		return nil, fmt.Errorf("error getting order (%s) details: response empty", orderId)
	}

	rst := make([]string, len(resp.Resources))
	for i, v := range resp.Resources {
		rst[i] = v.ResourceId
	}
	return rst, nil
}
