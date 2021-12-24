package products

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

// Get products
func Get(client *golangsdk.ServiceClient, engine string) (*GetResponse, error) {
	if len(engine) == 0 {
		return nil, fmt.Errorf("The parameter \"engine\" cannot be empty, it is required.")
	}
	url := getURL(client)
	url = url + "?engine=" + engine

	var rst golangsdk.Result
	_, err := client.Get(url, &rst.Body, nil)
	if err == nil {
		var r GetResponse
		err = rst.ExtractInto(&r)
		return &r, err
	}
	return nil, err
}
