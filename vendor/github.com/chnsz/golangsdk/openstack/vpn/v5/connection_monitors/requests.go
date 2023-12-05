package connection_monitors

import (
	"github.com/chnsz/golangsdk"
)

func List(client *golangsdk.ServiceClient) (r CommonResult) {
	_, r.Err = client.Get(listURL(client), &r.Body, nil)
	return
}
