package vpcs

import (
	"reflect"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.

type ListOpts struct {
	// ID is the unique identifier for the vpc.
	ID string `json:"id"`

	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	Name string `json:"name"`

	//Specifies the range of available subnets in the VPC.
	CIDR string `json:"cidr"`

	// Enterprise project ID.
	EnterpriseProjectID string `q:"enterprise_project_id"`

	// Status indicates whether or not a vpc is currently operational.
	Status string `json:"status"`

	// Specifies tags VPCs must match (returning those matching all tags).
	Tags string `q:"tags"`

	// Specifies tags VPCs must match (returning those matching at least one of the tags).
	TagsAny string `q:"tags-any"`

	// Specifies tags VPCs mustn't match (returning those missing all tags).
	NotTags string `q:"not-tags"`

	// Specifies tags VPCs mustn't match (returning those missing at least one of the tags).
	NotTagsAny string `q:"not-tags-any"`
}

func (opts ListOpts) hasQueryParameter() bool {
	return opts.EnterpriseProjectID != "" || opts.Tags != "" || opts.TagsAny != "" || opts.NotTags != "" || opts.NotTagsAny != ""
}

// ToVpcListQuery formats a ListOpts into a query string
func (opts ListOpts) ToVpcListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns collection of
// vpcs. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those vpcs that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Vpc, error) {
	url := rootURL(c)
	if opts.hasQueryParameter() {
		query, err := opts.ToVpcListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VpcPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()
	if err != nil {
		return nil, err
	}

	allVpcs, err := ExtractVpcs(pages)
	if err != nil {
		return nil, err
	}

	return FilterVPCs(allVpcs, opts)
}

func FilterVPCs(vpcs []Vpc, opts ListOpts) ([]Vpc, error) {

	var refinedVPCs []Vpc
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if opts.CIDR != "" {
		m["CIDR"] = opts.CIDR
	}

	if len(m) > 0 && len(vpcs) > 0 {
		for _, vpc := range vpcs {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&vpc, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedVPCs = append(refinedVPCs, vpc)
			}
		}

	} else {
		refinedVPCs = vpcs
	}

	return refinedVPCs, nil
}

func getStructField(v *Vpc, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVpcCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new vpc. There are
// no required values.
type CreateOpts struct {
	Name                       string `json:"name,omitempty"`
	CIDR                       string `json:"cidr,omitempty"`
	Description                string `json:"description,omitempty"`
	EnhancedLocalRoute         *bool  `json:"enhanced_local_route,omitempty"`
	EnterpriseProjectID        string `json:"enterprise_project_id,omitempty"`
	BlockServiceEndpointStates string `json:"block_service_endpoint_states,omitempty"`
}

// ToVpcCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToVpcCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "vpc")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical vpc. When it is created, the vpc does not have an internal
// interface - it is not associated to any subnet.
//
// You can optionally specify an external gateway for a vpc using the
// GatewayInfo struct. The external gateway for the vpc must be plugged into
// an external network (it is external if its `vpc:external' field is set to
// true).
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVpcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular vpc based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVpcUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a vpc.
type UpdateOpts struct {
	Name               string  `json:"name,omitempty"`
	CIDR               string  `json:"cidr,omitempty"`
	Description        *string `json:"description,omitempty"`
	EnhancedLocalRoute *bool   `json:"enhanced_local_route,omitempty"`
	EnableSharedSnat   *bool   `json:"enable_shared_snat,omitempty"`
}

// ToVpcUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToVpcUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "vpc")
}

// Update allows vpcs to be updated. You can update the name, administrative
// state, and the external gateway. For more information about how to set the
// external gateway for a vpc, see Create. This operation does not enable
// the update of vpc interfaces. To do this, use the AddInterface and
// RemoveInterface functions.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular vpc based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
