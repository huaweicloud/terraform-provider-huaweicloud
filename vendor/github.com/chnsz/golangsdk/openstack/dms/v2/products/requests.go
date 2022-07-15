package products

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

// Get products
//
// Deprecated: Please use the GetFlavor method to query product (new format) details of DMS kafka.
// The old format is '00300-30308-0--0', but the new that we recommanded is 'c6.2u4g.cluster'.
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

type ListOpts struct {
	// The product ID.
	ProductId string `q:"product_id"`
}

// List is a method to query the list of product details using given parameters.
func List(c *golangsdk.ServiceClient, engineType string, opts ListOpts) (*ListResp, error) {
	url := listURL(c, engineType)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r ListResp
	_, err = c.Get(url, &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
