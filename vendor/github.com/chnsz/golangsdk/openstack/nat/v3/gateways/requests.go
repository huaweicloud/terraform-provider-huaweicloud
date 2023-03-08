package gateways

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateOpts is the structure used to create a new private NAT gateway.
type CreateOpts struct {
	// The name of the private NAT gateway.
	// The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
	Name string `json:"name" required:"true"`
	// The subnet configuration of the private NAT gateway.
	DownLinkVpcs []DownLinkVpc `json:"downlink_vpcs" required:"true"`
	// The description of the private NAT gateway, which contain maximum of `255` characters, and
	// angle brackets (<) and (>) are not allowed.
	Description string `json:"description,omitempty"`
	// The specification of the private NAT gateway.
	// The valid values are as follows:
	// + **Small**: Small type, which supports up to `20` rules, `200 Mbit/s` bandwidth, `20,000` PPS and `2,000` SNAT
	//   connections.
	// + **Medium**: Medium type, which supports up to `50` rules, `500 Mbit/s` bandwidth, `50,000` PPS and `5,000` SNAT
	//   connections.
	// + **Large**: Large type, which supports up to `200` rules, `2 Gbit/s` bandwidth, `200,000` PPS and `20,000` SNAT
	//   connections.
	// + **Extra-Large**: Extra-large type, which supports up to `500` rules, `5 Gbit/s` bandwidth, `500,000` PPS and
	//   `50,000` SNAT connections.
	Spec string `json:"spec,omitempty"`
	// The enterprise project ID to which the private NAT gateway belongs.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// The key/value pairs to associate with the NAT geteway.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// DownLinkVpc is an object that represents the subnet configuration to which private NAT gateway belongs.
type DownLinkVpc struct {
	SubnetId string `json:"virsubnet_id" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new private NAT gateway using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Gateway, error) {
	b, err := golangsdk.BuildRequestBody(opts, "gateway")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Gateway, err
}

// Get is a method used to obtain the private NAT gateway detail by its ID.
func Get(c *golangsdk.ServiceClient, gatewayId string) (*Gateway, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, gatewayId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Gateway, err
}

// UpdateOpts is the structure used to modify an existing private NAT gateway.
type UpdateOpts struct {
	// The name of the private NAT gateway.
	// The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
	Name string `json:"name,omitempty"`
	// The description of the private NAT gateway, which contain maximum of `255` characters, and
	// angle brackets (<) and (>) are not allowed.
	Description *string `json:"description,omitempty"`
	// The specification of the private NAT gateway.
	// The valid values are as follows:
	// + **Small**: Small type, which supports up to `20` rules, `200 Mbit/s` bandwidth, `20,000` PPS and `2,000` SNAT
	//   connections.
	// + **Medium**: Medium type, which supports up to `50` rules, `500 Mbit/s` bandwidth, `50,000` PPS and `5,000` SNAT
	//   connections.
	// + **Large**: Large type, which supports up to `200` rules, `2 Gbit/s` bandwidth, `200,000` PPS and `20,000` SNAT
	//   connections.
	// + **Extra-Large**: Extra-large type, which supports up to `500` rules, `5 Gbit/s` bandwidth, `500,000` PPS and
	//   `50,000` SNAT connections.
	Spec string `json:"spec,omitempty"`
}

// Update is a method used to modify an existing private NAT gateway using given parameters.
func Update(c *golangsdk.ServiceClient, gatewayId string, opts UpdateOpts) (*Gateway, error) {
	b, err := golangsdk.BuildRequestBody(opts, "gateway")
	if err != nil {
		return nil, err
	}

	var r updateResp
	_, err = c.Put(resourceURL(c, gatewayId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Gateway, err
}

// Delete is a method to remove the specified private NAT gateway using its ID.
func Delete(c *golangsdk.ServiceClient, gatewayId string) error {
	_, err := c.Delete(resourceURL(c, gatewayId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
