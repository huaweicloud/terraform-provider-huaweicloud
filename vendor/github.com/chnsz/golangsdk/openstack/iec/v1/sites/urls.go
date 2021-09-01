package sites

import (
	"github.com/chnsz/golangsdk"
)

func ListURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("sites")
}
