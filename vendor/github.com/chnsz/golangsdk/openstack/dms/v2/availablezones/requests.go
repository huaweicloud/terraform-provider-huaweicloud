package availablezones

import (
	"github.com/chnsz/golangsdk"
)

// Get available zones
func Get(client *golangsdk.ServiceClient) (*GetResponse, error) {
	var rst golangsdk.Result
	_, err := client.Get(getURL(client), &rst.Body, nil)
	if err == nil {
		var r GetResponse
		err = rst.ExtractInto(&r)
		return &r, err
	}
	return nil, err
}
