package resolverrule

import "github.com/chnsz/golangsdk"

func baseURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("resolverrules")
}

func resourceURL(client *golangsdk.ServiceClient, resolverRuleID string) string {
	return client.ServiceURL("resolverrules", resolverRuleID)
}
