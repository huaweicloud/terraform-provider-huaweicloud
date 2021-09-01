package images

import "github.com/chnsz/golangsdk"

// ListURL list iec image url
func ListURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("images")
}
