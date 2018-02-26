package healthcheck

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToHealthCheckCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	ListenerID             string `json:"listener_id" required:"true"`
	HealthcheckProtocol    string `json:"healthcheck_protocol,omitempty"`
	HealthcheckUri         string `json:"healthcheck_uri,omitempty"`
	HealthcheckConnectPort int    `json:"healthcheck_connect_port,omitempty"`
	HealthyThreshold       int    `json:"healthy_threshold,omitempty"`
	UnhealthyThreshold     int    `json:"unhealthy_threshold,omitempty"`
	HealthcheckTimeout     int    `json:"healthcheck_timeout,omitempty"`
	HealthcheckInterval    int    `json:"healthcheck_interval,omitempty"`
}

// ToHealthCheckCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToHealthCheckCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToHealthCheckCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create ELB-HealthCheck url:%q, body=%#v", rootURL(c), b)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToHealthCheckUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	HealthcheckProtocol    string `json:"healthcheck_protocol,omitempty"`
	HealthcheckUri         string `json:"healthcheck_uri,omitempty"`
	HealthcheckConnectPort int    `json:"healthcheck_connect_port,omitempty"`
	HealthyThreshold       int    `json:"healthy_threshold,omitempty"`
	UnhealthyThreshold     int    `json:"unhealthy_threshold,omitempty"`
	HealthcheckTimeout     int    `json:"healthcheck_timeout,omitempty"`
	HealthcheckInterval    int    `json:"healthcheck_interval,omitempty"`
}

func (u UpdateOpts) IsNeedUpdate() (bool, error) {
	d, e := u.ToHealthCheckUpdateMap()
	if e == nil {
		return len(d) != 0, nil
	}
	return false, e
}

// ToHealthCheckUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToHealthCheckUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified HealthCheck.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToHealthCheckUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular HealthCheck based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
