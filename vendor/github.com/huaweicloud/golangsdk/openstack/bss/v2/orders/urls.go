package orders

import "github.com/huaweicloud/golangsdk"

func unsubscribeURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("orders/subscriptions/resources/unsubscribe")
}
