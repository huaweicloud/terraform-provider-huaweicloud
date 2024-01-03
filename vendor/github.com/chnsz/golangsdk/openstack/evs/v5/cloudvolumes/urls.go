package cloudvolumes

import "github.com/chnsz/golangsdk"

func qoSURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("cloudvolumes", id, "qos")
}
