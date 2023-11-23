package endpoints

import (
	"encoding/json"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ListOpts struct {
	Direction string `q:"direction" required:"true"`
	VpcID     string `q:"vpc_id,omitempty"`
	Limit     int    `q:"limit,omitempty"`
	Offset    int    `q:"offset,omitempty"`
}

type CreateOpt struct {
	Name        string        `json:"name" required:"true"`
	Direction   string        `json:"direction" required:"true"`
	Region      string        `json:"region" required:"true"`
	IPAddresses []IPAddresses `json:"ipaddresses" required:"true"`
}

type IPAddresses struct {
	SubnetID string `json:"subnet_id" required:"true"`
	IP       string `json:"ip,omitempty"`
}

type UpdateOpts struct {
	Name string `json:"name" required:"true"`
}

func Create(c *golangsdk.ServiceClient, opts CreateOpt) (r CreateResult) {
	_, err := c.Post(baseUrl(c), opts, &r.Body, nil)
	if err != nil {
		r.Err = err
	}
	return r
}

func Get(c *golangsdk.ServiceClient, endpointID string) (r GetResult) {
	_, r.Err = c.Get(resourceUrl(c, endpointID), &r.Body, nil)
	return r
}

func Update(c *golangsdk.ServiceClient, endpointID string, opts UpdateOpts) (r UpdateResult) {
	_, r.Err = c.Put(resourceUrl(c, endpointID), opts, &r.Body, nil)
	return
}

func Delete(c *golangsdk.ServiceClient, endpointID string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceUrl(c, endpointID), nil)
	return
}

func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Endpoint, error) {
	url := baseUrl(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	resp, err := pagination.ListAllItems(c, pagination.Offset, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var e struct {
		Instances []Endpoint `json:"endpoints"`
	}
	if err = json.Unmarshal(body, &e); err != nil {
		return nil, err
	}
	return e.Instances, nil
}
