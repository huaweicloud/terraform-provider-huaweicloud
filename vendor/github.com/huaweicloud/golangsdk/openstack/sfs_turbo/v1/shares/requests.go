package shares

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create an SFS Turbo. This object is
// passed to shares.Create().
type CreateOpts struct {
	// Defines the SFS Turbo file system name
	Name string `json:"name" required:"true"`
	// Defines the SFS Turbo file system protocol to use, the vaild value is NFS.
	ShareProto string `json:"share_proto,omitempty"`
	// ShareType defines the file system type. the vaild values are STANDARD and PERFORMANCE.
	ShareType string `json:"share_type" required:"true"`
	// Size in GB, range from 500 to 32768.
	Size int `json:"size" required:"true"`
	// The availability zone of the SFS Turbo file system
	AvailabilityZone string `json:"availability_zone" required:"true"`
	// The VPC ID
	VpcID string `json:"vpc_id" required:"true"`
	// The subnet ID
	SubnetID string `json:"subnet_id" required:"true"`
	// The security group ID
	SecurityGroupID string `json:"security_group_id" required:"true"`
	// The enterprise project ID
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// The backup ID
	BackupID string `json:"backup_id,omitempty"`
	// Share description
	Description string `json:"description,omitempty"`
	// The metadata information
	Metadata Metadata `json:"metadata,omitempty"`
}

// Metadata specifies the metadata information
type Metadata struct {
	ExpandType            string `json:"expand_type,omitempty"`
	CryptKeyID            string `json:"crypt_key_id,omitempty"`
	DedicatedFlavor       string `json:"dedicated_flavor,omitempty"`
	MasterDedicatedHostID string `json:"master_dedicated_host_id,omitempty"`
	SlaveDedicatedHostID  string `json:"slave_dedicated_host_id,omitempty"`
	DedicatedStorageID    string `json:"dedicated_storage_id,omitempty"`
}

// ToShareCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "share")
}

// Create will create a new SFS Turbo file system based on the values in CreateOpts. To extract
// the Share object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// List returns a Pager which allows you to iterate over a collection of
// SFS Turbo resources.
func List(c *golangsdk.ServiceClient) ([]Turbo, error) {
	pages, err := pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return TurboPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()
	if err != nil {
		return nil, err
	}

	return ExtractTurbos(pages)
}

// Get will get a single SFS Trubo file system with given UUID
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)

	return
}

// Delete will delete an existing SFS Trubo file system with the given UUID.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

// ExpandOptsBuilder allows extensions to add additional parameters to the
// Expand request.
type ExpandOptsBuilder interface {
	ToShareExpandMap() (map[string]interface{}, error)
}

// ExpandOpts contains the options for expanding a SFS Turbo. This object is
// passed to shares.Expand().
type ExpandOpts struct {
	// Specifies the extend object.
	Extend ExtendOpts `json:"extend" required:"true"`
}

type ExtendOpts struct {
	// Specifies the post-expansion capacity (GB) of the shared file system.
	NewSize int `json:"new_size" required:"true"`
}

// ToShareExpandMap assembles a request body based on the contents of a
// ExpandOpts.
func (opts ExpandOpts) ToShareExpandMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Expand will expand a SFS Turbo based on the values in ExpandOpts.
func Expand(client *golangsdk.ServiceClient, share_id string, opts ExpandOptsBuilder) (r ExpandResult) {
	b, err := opts.ToShareExpandMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, share_id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}
