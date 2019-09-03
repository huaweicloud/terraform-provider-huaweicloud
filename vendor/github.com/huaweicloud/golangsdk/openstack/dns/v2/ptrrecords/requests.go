package ptrrecords

import (
	"github.com/huaweicloud/golangsdk"
)

// Get returns information about a ptr, given its ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional attributes to the
// Create request.
type CreateOptsBuilder interface {
	ToPtrCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies the attributes used to create a ptr.
type CreateOpts struct {
	// Name of the ptr.
	PtrName string `json:"ptrdname" required:"true"`

	// Description of the ptr.
	Description string `json:"description,omitempty"`

	// TTL is the time to live of the ptr.
	TTL int `json:"-"`

	// Tags of the ptr.
	Tags []Tag `json:"tags,omitempty"`
}

// Tag is a structure of key value pair.
type Tag struct {
	//tag key
	Key string `json:"key" required:"true"`
	//tag value
	Value string `json:"value" required:"true"`
}

// ToPtrCreateMap formats an CreateOpts structure into a request body.
func (opts CreateOpts) ToPtrCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.TTL > 0 {
		b["ttl"] = opts.TTL
	}

	return b, nil
}

// Create implements a ptr create/update request.
func Create(client *golangsdk.ServiceClient, region string, fip_id string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPtrCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(baseURL(client, region, fip_id), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete implements a ptr delete request.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	b := map[string]string{
		"ptrname": "null",
	}
	_, r.Err = client.Patch(resourceURL(client, id), &b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}
