package ipaddress

import "github.com/chnsz/golangsdk"

type CreateOpts struct {
	IPAddress `json:"ipaddress" required:"true"`
}

type IPAddress struct {
	SubnetID string `json:"subnet_id" required:"true"`
	IP       string `json:"ip,omitempty"`
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts, endpointID string) (r CreateResult) {
	body, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(baseUrl(c, endpointID), body, &r.Body, nil)
	return
}

func Delete(c *golangsdk.ServiceClient, endpointID string, ipaddressID string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceUrl(c, endpointID, ipaddressID), nil)
	return
}

func List(c *golangsdk.ServiceClient, endpointID string) (r ListResult) {
	_, r.Err = c.Get(baseUrl(c, endpointID), &r.Body, nil)
	return
}
