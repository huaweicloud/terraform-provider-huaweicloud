package products

import (
	"github.com/chnsz/golangsdk"
)

// Get products
func Get(client *golangsdk.ServiceClient, engine string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, engine), &r.Body, nil)
	return
}
