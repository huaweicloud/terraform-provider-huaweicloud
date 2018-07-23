package shares

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json",
		"X-Openstack-Manila-Api-Version": "2.9"},
}

// SortDir is a type for specifying in which direction to sort a list of Shares.
type SortDir string

// SortKey is a type for specifying by which key to sort a list of Shares.
type SortKey string

var (
	// SortAsc is used to sort a list of Shares in ascending order.
	SortAsc SortDir = "asc"
	// SortDesc is used to sort a list of Shares in descending order.
	SortDesc SortDir = "desc"
	// SortId is used to sort a list of Shares by id.
	SortId SortKey = "id"
	// SortName is used to sort a list of Shares by name.
	SortName SortKey = "name"
	// SortSize is used to sort a list of Shares by size.
	SortSize SortKey = "size"
	// SortHost is used to sort a list of Shares by host.
	SortHost SortKey = "host"
	// SortShareProto is used to sort a list of Shares by share_proto.
	SortShareProto SortKey = "share_proto"
	// SortStatus is used to sort a list of Shares by status.
	SortStatus SortKey = "status"
	// SortProjectId is used to sort a list of Shares by project_id.
	SortProjectId SortKey = "project_id"
	// SortShareTypeId is used to sort a list of Shares by share_type_id.
	SortShareTypeId SortKey = "share_type_id"
	// SortShareNetworkId is used to sort a list of Shares by share_network_id.
	SortShareNetworkId SortKey = "share_network_id"
	// SortSnapshotId is used to sort a list of Shares by snapshot_id.
	SortSnapshotId SortKey = "snapshot_id"
	// SortCreatedAt is used to sort a list of Shares by date created.
	SortCreatedAt SortKey = "created_at"
	// SortUpdatedAt is used to sort a list of Shares by date updated.
	SortUpdatedAt SortKey = "updated_at"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToShareListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the share attributes you want to see returned. SortKey allows you to sort
// by a particular share attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID       string
	Status   string  `q:"status"`
	Name     string  `q:"name"`
	Limit    int     `q:"limit"`
	Offset   int     `q:"offset"`
	SortKey  SortKey `q:"sort_key"`
	SortDir  SortDir `q:"sort_dir"`
	IsPublic bool    `q:"is_public"`
}

// List returns a Pager which allows you to iterate over a collection of
// share resources. It accepts a ListOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Share, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return SharePage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allShares, err := ExtractShares(pages)
	if err != nil {
		return nil, err
	}

	return FilterShares(allShares, opts)
}

func FilterShares(shares []Share, opts ListOpts) ([]Share, error) {

	var refinedShares []Share
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}

	if len(m) > 0 && len(shares) > 0 {
		for _, share := range shares {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&share, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedShares = append(refinedShares, share)
			}
		}
	} else {
		refinedShares = shares
	}
	return refinedShares, nil
}

func getStructField(v *Share, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create a Share. This object is
// passed to shares.Create(). For more information about these parameters,
// please refer to the Share object, or the shared file systems API v2
// documentation
type CreateOpts struct {
	// Defines the share protocol to use
	ShareProto string `json:"share_proto" required:"true"`
	// Size in GB
	Size int `json:"size" required:"true"`
	// Defines the share name
	Name string `json:"name,omitempty"`
	// Share description
	Description string `json:"description,omitempty"`
	// ShareType defines the sharetype. If omitted, a default share type is used
	ShareType string `json:"share_type,omitempty"`
	// The UUID from which to create a share
	SnapshotID string `json:"snapshot_id,omitempty"`
	// Determines whether or not the share is public
	IsPublic bool `json:"is_public,omitempty"`
	// Key value pairs of user defined metadata
	Metadata map[string]string `json:"metadata,omitempty"`
	// The UUID of the share network to which the share belongs to
	ShareNetworkID string `json:"share_network_id,omitempty"`
	// The UUID of the consistency group to which the share belongs to
	ConsistencyGroupID string `json:"consistency_group_id,omitempty"`
	// The availability zone of the share
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

// ToShareCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "share")
}

// Create will create a new Share based on the values in CreateOpts. To extract
// the Share object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToShareUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a Share.
type UpdateOpts struct {
	// DisplayName is equivalent to Name. The API supports using both
	// This is an inherited attribute from the block storage API
	DisplayName string `json:"display_name" required:"true"`
	// DisplayDescription is equivalent to Description. The API supports using bot
	// This is an inherited attribute from the block storage API
	DisplayDescription string `json:"display_description,omitempty"`
}

// ToShareUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToShareUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "share")
}

// Update allows shares to be updated. You can update the DisplayName, DisplayDescription.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToShareUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get will get a single share with given UUID
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)

	return
}

// Delete will delete an existing Share with the given UUID.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

// ListAccessRights lists all access rules assigned to a Share based on its id. To extract
// the AccessRight slice from the response, call the Extract method on the AccessRightsResult.
// Client must have Microversion set; minimum supported microversion for ListAccessRights is 2.7.
func ListAccessRights(client *golangsdk.ServiceClient, share_id string) (r AccessRightsResult) {
	requestBody := map[string]interface{}{"os-access_list": nil}
	_, r.Err = client.Post(rootURL(client, share_id), requestBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// GrantAccessOptsBuilder allows extensions to add additional parameters to the
// GrantAccess request.
type GrantAccessOptsBuilder interface {
	ToGrantAccessMap() (map[string]interface{}, error)
}

// GrantAccessOpts contains the options for creation of an GrantAccess request.
// For more information about these parameters, please, refer to the shared file systems API v2,
// Share Actions, Grant Access documentation
type GrantAccessOpts struct {
	// The access rule type that can be "ip", "cert" or "user".
	AccessType string `json:"access_type"`
	// The value that defines the access that can be a valid format of IP, cert or user.
	AccessTo string `json:"access_to"`
	// The access level to the share is either "rw" or "ro".
	AccessLevel string `json:"access_level"`
}

// ToGrantAccessMap assembles a request body based on the contents of a
// GrantAccessOpts.
func (opts GrantAccessOpts) ToGrantAccessMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "os-allow_access")
}

// GrantAccess will grant access to a Share based on the values in GrantAccessOpts. To extract
// the GrantAccess object from the response, call the Extract method on the GrantAccessResult.
// Client must have Microversion set; minimum supported microversion for GrantAccess is 2.7.
func GrantAccess(client *golangsdk.ServiceClient, share_id string, opts GrantAccessOptsBuilder) (r GrantAccessResult) {
	b, err := opts.ToGrantAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, share_id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete the Access Rule
type DeleteAccessOptsBuilder interface {
	ToDeleteAccessMap() (map[string]interface{}, error)
}

type DeleteAccessOpts struct {
	// The access ID to be deleted
	AccessID string `json:"access_id"`
}

func (opts DeleteAccessOpts) ToDeleteAccessMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "os-deny_access")
}

//Deletes the Access Rule
func DeleteAccess(client *golangsdk.ServiceClient, share_id string, opts DeleteAccessOptsBuilder) (r DeleteAccessResult) {
	b, err := opts.ToDeleteAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, share_id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

//Gets the Mount/Export Locations of the SFS specified
func GetExportLocations(client *golangsdk.ServiceClient, id string) (r GetExportLocationsResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = client.Get(getMountLocationsURL(client, id), &r.Body, reqOpt)
	return
}

// ExpandOptsBuilder allows extensions to add additional parameters to the
// Expand request.
type ExpandOptsBuilder interface {
	ToShareExpandMap() (map[string]interface{}, error)
}

// ExpandOpts contains the options for expanding a Share. This object is
// passed to shares.Expand(). For more information about these parameters,
// please refer to the Share object, or the shared file systems API v2
// documentation
type ExpandOpts struct {
	// Specifies the os-extend object.
	OSExtend OSExtendOpts `json:"os-extend" required:"true"`
}

type OSExtendOpts struct {
	// Specifies the post-expansion capacity (GB) of the shared file system.
	NewSize int `json:"new_size" required:"true"`
}

// ToShareExpandMap assembles a request body based on the contents of a
// ExpandOpts.
func (opts ExpandOpts) ToShareExpandMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Expand will expand a  Share based on the values in ExpandOpts.
func Expand(client *golangsdk.ServiceClient, share_id string, opts ExpandOptsBuilder) (r ExpandResult) {
	b, err := opts.ToShareExpandMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(grantAccessURL(client, share_id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// ShrinkOptsBuilder allows extensions to add additional parameters to the
// Shrink request.
type ShrinkOptsBuilder interface {
	ToShareShrinkMap() (map[string]interface{}, error)
}

// ShrinkOpts contains the options for shrinking a Share. This object is
// passed to shares.Shrink(). For more information about these parameters,
// please refer to the Share object, or the shared file systems API v2
// documentation
type ShrinkOpts struct {
	// Specifies the os-shrink object.
	OSShrink OSShrinkOpts `json:"os-shrink" required:"true"`
}

type OSShrinkOpts struct {
	// Specifies the post-shrinking capacity (GB) of the shared file system.
	NewSize int `json:"new_size" required:"true"`
}

// ToShareShrinkMap assembles a request body based on the contents of a
// ShrinkOpts.
func (opts ShrinkOpts) ToShareShrinkMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Shrink will shrink a  Share based on the values in ShrinkOpts.
func Shrink(client *golangsdk.ServiceClient, share_id string, opts ShrinkOptsBuilder) (r ShrinkResult) {
	b, err := opts.ToShareShrinkMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(grantAccessURL(client, share_id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}
