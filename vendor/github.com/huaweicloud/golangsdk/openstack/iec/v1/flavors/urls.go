package flavors

import (
	"github.com/huaweicloud/golangsdk"
)

func ListURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudservers", "flavors")
}
