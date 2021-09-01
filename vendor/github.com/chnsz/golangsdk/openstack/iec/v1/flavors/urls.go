package flavors

import (
	"github.com/chnsz/golangsdk"
)

func ListURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudservers", "flavors")
}
