package orders

import "github.com/chnsz/golangsdk"

func unsubscribeURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("orders/subscriptions/resources/unsubscribe")
}

func getURL(sc *golangsdk.ServiceClient, id string) string {
	return sc.ServiceURL("orders/customer-orders/details", id)
}
