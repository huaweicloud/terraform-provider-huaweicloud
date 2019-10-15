package block_devices

import (
	"github.com/huaweicloud/golangsdk"
)

func Get(c *golangsdk.ServiceClient, server_id string, volume_id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, server_id, volume_id), &r.Body, nil)
	return
}
