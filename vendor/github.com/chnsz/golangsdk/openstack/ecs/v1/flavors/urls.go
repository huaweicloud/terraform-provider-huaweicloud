package flavors

import "github.com/chnsz/golangsdk"

func listURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("cloudservers", "flavors")
}
