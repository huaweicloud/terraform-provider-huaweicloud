package backendecs

import (
	"log"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToBackendECSCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	ServerId string `json:"server_id" required:"true"`
	Address  string `json:"address" required:"true"`
}

// ToBackendECSCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToBackendECSCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder, lId string) (r elb.JobResult) {
	b, err := opts.ToBackendECSCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	//API takes an array of these...
	body := []map[string]interface{}{b}
	log.Printf("[DEBUG] create ELB-BackendECS url:%q, body=%#v", rootURL(c, lId), body)

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, lId), body, &r.Body, reqOpt)
	return
}

type GetOptsBuilder interface {
	ToBackendECSListQuery() (string, error)
}

type getOpts struct {
	ID string `q:"id"`
}

func (opts getOpts) ToBackendECSListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *golangsdk.ServiceClient, lId string, backendId string) (r GetResult) {
	url := rootURL(c, lId)
	opts := getOpts{ID: backendId}
	query, err := opts.ToBackendECSListQuery()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] get ELB-BackendECS opt=%#v, query=%s", opts, query)
	url += query
	log.Printf("[DEBUG] get ELB-BackendECS url:%q, backendId:%q", url, backendId)

	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type DeleteOptsBuilder interface {
	ToBackendECSDeleteMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.

type RemoveMemberField struct {
	ID string `json:"id" required:"true"`
}

type DeleteOpts struct {
	RemoveMember []RemoveMemberField `json:"removeMember" required:"true"`
}

// ToBackendECSUpdateMap casts a UpdateOpts struct to a map.
func (opts DeleteOpts) ToBackendECSDeleteMap() (map[string]interface{}, error) {
	/*
		o, err := golangsdk.BuildRequestBody(opts, "")
		if err != nil {
			return nil, err
		}
		rm = removeMeber{RemoveMember: []map[string]string{o.([string]string)}}
	*/
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified BackendECS.
func Delete(c *golangsdk.ServiceClient, lId string, opts DeleteOpts) (r elb.JobResult) {
	b, err := opts.ToBackendECSDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] deleting ELB-BackendECS, request opts=%#v", b)

	_, r.Err = c.Post(actionURL(c, lId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
