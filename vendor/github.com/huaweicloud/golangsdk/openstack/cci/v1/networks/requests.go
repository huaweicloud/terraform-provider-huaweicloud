package networks

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToNetworkCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new network
type CreateOpts struct {
	// API type, fixed value Network
	Kind string `json:"kind" required:"true"`
	// API version, fixed value networking.cci.io
	ApiVersion string `json:"apiVersion" required:"true"`
	// Metadata required to create a network
	Metadata CreateMetaData `json:"metadata" required:"true"`
	// Specifications to create a network
	Spec Spec `json:"spec" required:"true"`
}

// Metadata required to create a network
type CreateMetaData struct {
	//Network unique name
	Name string `json:"name" required:"true"`
	//Network annotation, key/value pair format
	Annotations map[string]string `json:"annotations" required:"true"`
}

// Specifications to create a network
type Spec struct {
	// Network CIDR
	Cidr string `json:"type,omitempty"`
	// Network VPC ID
	AttachedVPC string `json:"attachedVPC" required:"true"`
	// Network Type
	NetworkType string `json:"networkType" required:"true"`
	// Network ID
	NetworkID string `json:"networkID" required:"true"`
	// Subnet ID
	SubnetID string `json:"subnetID" required:"true"`
	// Network AZ
	AvailableZone string `json:"availableZone" required:"true"`
}

// ToNetworkCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToNetworkCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new network.
func Create(c *golangsdk.ServiceClient, ns string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNetworkCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c, ns), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular network based on its unique ID.
func Get(c *golangsdk.ServiceClient, ns, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, ns, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// Delete will permanently delete a particular network based on its unique ID.
func Delete(c *golangsdk.ServiceClient, ns, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, ns, id), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
