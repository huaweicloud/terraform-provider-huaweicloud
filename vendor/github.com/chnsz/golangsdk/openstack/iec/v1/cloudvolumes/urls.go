package cloudvolumes

import (
	"github.com/chnsz/golangsdk"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudvolumes")
}

func GetURL(c *golangsdk.ServiceClient, CloudVolumeID string) string {
	return c.ServiceURL("cloudvolumes", CloudVolumeID)
}

func ListVolumeTypeURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudvolumes", "volume-types")
}
