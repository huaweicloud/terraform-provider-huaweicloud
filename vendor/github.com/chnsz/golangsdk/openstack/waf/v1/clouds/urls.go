package clouds

import "github.com/chnsz/golangsdk"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("subscription/purchase/prepaid-cloud-waf")
}

func createOrDeletePostPaidURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("postpaid")
}

func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("subscription")
}

func updateURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("subscription/batchalter/prepaid-cloud-waf")
}
