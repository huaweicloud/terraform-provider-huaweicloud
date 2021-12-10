package bandwidths

import (
	"github.com/chnsz/golangsdk"
)

func GetURL(c *golangsdk.ServiceClient, bandwidthID string) string {
	return c.ServiceURL("bandwidths", bandwidthID)
}

func UpdateURL(c *golangsdk.ServiceClient, bandwidthID string) string {
	return c.ServiceURL("bandwidths", bandwidthID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("bandwidths")
}
