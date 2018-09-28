package maintainwindows

import (
	"github.com/huaweicloud/golangsdk"
)

// Get maintain windows
func Get(client *golangsdk.ServiceClient) (r GetResult) {
	_, r.Err = client.Get(getURL(client), &r.Body, nil)
	return
}
