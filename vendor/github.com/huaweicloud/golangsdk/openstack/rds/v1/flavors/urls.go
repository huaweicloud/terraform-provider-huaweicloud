package flavors

import "github.com/huaweicloud/golangsdk"

func listURL(c *golangsdk.ServiceClient, dataStoreID string, region string) string {

	return c.ResourceBaseURL() + "flavors?dbId=" + dataStoreID + "&region=" + region
}
