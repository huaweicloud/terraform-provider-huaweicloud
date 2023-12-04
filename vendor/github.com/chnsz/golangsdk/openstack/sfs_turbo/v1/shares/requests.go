package shares

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the structure used to create a new SFS Turbo resource.
type CreateOpts struct {
	// Turbo configuration details.
	Share Share `json:"share" required:"share"`
	// The configuration of pre-paid billing mode.
	BssParam *BssParam `json:"bss_param,omitempty"`
}

// CreateOpts contains the options for create an SFS Turbo. This object is
// passed to shares.Create().
type Share struct {
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
	HpcBw                 string `json:"hpc_bw,omitempty"`
}

// BssParam is an object that represents the prepaid configuration.
type BssParam struct {
	// The number of cycles for prepaid.
	// + minimum: 1
	// + maximum: 11
	PeriodNum int `json:"period_num" required:"true"`
	// The prepaid type.
	// + 2: month
	// + 3: year
	PeriodType int `json:"period_type" require:"true"`
	// Whether to automatically renew.
	// + 0: manual renew.
	// + 1: automatic renew.
	IsAutoRenew *int `json:"is_auto_renew,omitempty"`
	// Whether to pay automatically.
	// + 0: manual payment.
	// + 1: automatic payment.
	IsAutoPay *int `json:"is_auto_pay,omitempty"`
}

// ToShareCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
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

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToVolumeListQuery() (string, error)
}

// ListOpts holds options for listing Volumes. It is passed to the volumes.List
// function.
type ListOpts struct {
	// Requests a page size of items.
	Limit int `q:"limit"`
	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`
}

// ToVolumeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// ListPage returns one page limited by the conditions provided in Opts.
func ListPage(client *golangsdk.ServiceClient, opts ListOptsBuilder) (*PagedList, error) {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToVolumeListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	var rst golangsdk.Result
	_, err := client.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err != nil {
		return nil, err
	}

	var r PagedList
	err = rst.ExtractInto(&r)
	return &r, err
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

type UpdateNameOptsBuilder interface {
	ToShareUpdateNameMap() (map[string]interface{}, error)
}

type UpdateSecurityGroupIdOptsBuilder interface {
	ToShareUpdateSecurityGroupIdMap() (map[string]interface{}, error)
}

// ExpandOpts contains the options for expanding a SFS Turbo. This object is
// passed to shares.Expand().
type ExpandOpts struct {
	// Specifies the extend object.
	Extend ExtendOpts `json:"extend" required:"true"`
}

type UpdateNameOpts struct {
	Name string `json:"name" required:"true"`
}

type UpdateSecurityGroupIdOpts struct {
	SecurityGroupId string `json:"security_group_id" required:"true"`
}

// BssParamExtend is an object that represents the payment detail.
type BssParamExtend struct {
	// Whether to pay automatically.
	// + 0: manual payment.
	// + 1: automatic payment.
	IsAutoPay *int `json:"is_auto_pay,omitempty"`
}

type ExtendOpts struct {
	// Specifies the post-expansion capacity (GB) of the shared file system.
	NewSize int `json:"new_size" required:"true"`
	// New bandwidth of the file system, in GB/s. This parameter is only supported for HPC Cache file systems.
	NewBandwidth int `json:"new_bandwidth,omitempty"`
	// The configuration of pre-paid billing mode.
	BssParam *BssParamExtend `json:"bss_param,omitempty"`
}

// ToShareExpandMap assembles a request body based on the contents of a
// ExpandOpts.
func (opts ExpandOpts) ToShareExpandMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts UpdateNameOpts) ToShareUpdateNameMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "change_name")
}

func (opts UpdateSecurityGroupIdOpts) ToShareUpdateSecurityGroupIdMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "change_security_group")
}

// Expand will expand a SFS Turbo based on the values in ExpandOpts.
func Expand(client *golangsdk.ServiceClient, shareId string, opts ExpandOptsBuilder) (r ExpandResult) {
	b, err := opts.ToShareExpandMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, shareId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

func UpdateName(client *golangsdk.ServiceClient, shareId string, opts UpdateNameOptsBuilder) (r UpdateNameResult) {
	b, err := opts.ToShareUpdateNameMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, shareId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

func UpdateSecurityGroupId(client *golangsdk.ServiceClient, shareId string, opts UpdateSecurityGroupIdOptsBuilder) (r UpdateSecurityGroupIdResult) {
	b, err := opts.ToShareUpdateSecurityGroupIdMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, shareId), b, &r.Body, &golangsdk.RequestOpts{})
	return
}
