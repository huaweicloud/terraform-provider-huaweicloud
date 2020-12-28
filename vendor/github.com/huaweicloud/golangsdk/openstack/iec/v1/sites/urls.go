package sites

import (
	"github.com/huaweicloud/golangsdk"
)

func ListURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("sites")
}
