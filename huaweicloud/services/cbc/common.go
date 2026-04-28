package cbc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/resources"
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

// GetAutoPay is a method to return whether order is auto pay according to the user input.
// auto_pay parameter inputs and returns:
//
//	false: false
//	true, empty: true
//
// Before using this function, make sure the parameter behavior is auto pay (the default value is "true").
func GetAutoPay(d *schema.ResourceData) string {
	if val, ok := d.GetOk("auto_pay"); ok && val.(string) == "false" {
		return "false"
	}
	return "true"
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
		return nil, errors.New("order id is empty")
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
