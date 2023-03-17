package transitips

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateOpts is the structure used to create a new transit IP.
type CreateOpts struct {
	// The ID of the subnet to which the transit IP belongs.
	SubnetId string `json:"virsubnet_id" required:"true"`
	// The IP address
	IpAddress string `json:"ip_address,omitempty"`
	// The ID of the enterprise project to which the transit IP belongs.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// The key/value pairs to associate with the transit IP.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new transit IP using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*TransitIp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "transit_ip")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.TransitIp, err
}

// Get is a method used to obtain the transit IP detail by its ID.
func Get(c *golangsdk.ServiceClient, transitIpId string) (*TransitIp, error) {
	var r queryResp
	_, err := c.Get(resourceURL(c, transitIpId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.TransitIp, err
}

// Delete is a method to remove the specified transit IP using its ID.
func Delete(c *golangsdk.ServiceClient, transitIpId string) error {
	_, err := c.Delete(resourceURL(c, transitIpId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
