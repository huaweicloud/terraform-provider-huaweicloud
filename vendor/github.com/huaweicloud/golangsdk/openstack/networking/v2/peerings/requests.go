package peerings

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering  of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned.
type ListOpts struct {
	//ID is the unique identifier for the vpc_peering_connection.
	ID string `q:"id"`

	//Name is the human readable name for the vpc_peering_connection. It does not have to be
	// unique.
	Name string `q:"name"`

	//Status indicates whether or not a vpc_peering_connection is currently operational.
	Status string `q:"status"`

	// TenantId indicates  vpc_peering_connection avalable in specific tenant.
	TenantId string `q:"tenant_id"`

	// VpcId indicates vpc_peering_connection avalable in specific vpc.
	VpcId string `q:"vpc_id"`

	// VpcId indicates vpc_peering_connection available in specific vpc.
	Peer_VpcId string
}

func FilterVpcIdParam(opts ListOpts) (filter ListOpts) {

	if opts.VpcId != "" {
		filter.VpcId = opts.VpcId
	} else {
		filter.VpcId = opts.Peer_VpcId
	}

	filter.Name = opts.Name
	filter.ID = opts.ID
	filter.Status = opts.Status
	filter.TenantId = opts.TenantId

	return filter
}

// List returns a Pager which allows you to iterate over a collection of
// vpc_peering_connection  resources. It accepts a ListOpts struct, which allows you to
// filter  the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Peering, error) {
	filter := FilterVpcIdParam(opts)
	q, err := golangsdk.BuildQueryString(&filter)
	if err != nil {
		return nil, err
	}
	u := rootURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return PeeringConnectionPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allPeeringConns, err := ExtractPeerings(pages)
	if err != nil {
		return nil, err
	}

	return FilterVpcPeeringConns(allPeeringConns, opts)
}

func FilterVpcPeeringConns(peerings []Peering, opts ListOpts) ([]Peering, error) {
	var refinedPeerings []Peering
	var matched bool
	filterMap := map[string]interface{}{}

	if opts.VpcId != "" {
		filterMap["RequestVpcInfo"] = opts.VpcId
	}

	if opts.Peer_VpcId != "" {
		filterMap["AcceptVpcInfo"] = opts.Peer_VpcId
	}

	if len(filterMap) > 0 && len(peerings) > 0 {
		for _, peering := range peerings {
			matched = true

			for key, value := range filterMap {
				if sVal := getStructNestedField(&peering, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedPeerings = append(refinedPeerings, peering)
			}

		}
	} else {
		refinedPeerings = peerings
	}

	return refinedPeerings, nil

}

func getStructNestedField(v *Peering, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field).Interface()
	r1 := reflect.ValueOf(f)
	f1 := reflect.Indirect(r1).FieldByName("VpcId")
	return string(f1.String())
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Accept is used by a tenant to accept a VPC peering connection request initiated by another tenant.
func Accept(c *golangsdk.ServiceClient, id string) (r AcceptResult) {
	_, r.Err = c.Put(acceptURL(c, id), nil, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Reject is used by a tenant to reject a VPC peering connection request initiated by another tenant.
func Reject(c *golangsdk.ServiceClient, id string) (r RejectResult) {
	_, r.Err = c.Put(rejectURL(c, id), nil, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//CreateOptsBuilder is an interface by which can build the request body of vpc peering connection.
type CreateOptsBuilder interface {
	ToPeeringCreateMap() (map[string]interface{}, error)
}

//CreateOpts is a struct which is used to create vpc peering connection.
type CreateOpts struct {
	Name           string  `json:"name"`
	RequestVpcInfo VpcInfo `json:"request_vpc_info" required:"true"`
	AcceptVpcInfo  VpcInfo `json:"accept_vpc_info" required:"true"`
}

//ToVpcPeeringCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPeeringCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "peering")
}

//Create is a method by which can access to create the vpc peering connection.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPeeringCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

//Delete is a method by which can be able to delete a vpc peering connection.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

//UpdateOptsBuilder is an interface by which can be able to build the request body of vpc peering connection.
type UpdateOptsBuilder interface {
	ToVpcPeeringUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	Name string `json:"name,omitempty"`
}

//ToVpcPeeringUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToVpcPeeringUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "peering")
}

//Update is a method which can be able to update the name of vpc peering connection.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcPeeringUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
