package flavors

import "github.com/huaweicloud/golangsdk"

func listURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("cloudservers", "flavors")
}
