package eips

import "github.com/chnsz/golangsdk"

type GetResult struct {
	golangsdk.Result
}

// Get is a method by which can get the detailed information of public ip
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
