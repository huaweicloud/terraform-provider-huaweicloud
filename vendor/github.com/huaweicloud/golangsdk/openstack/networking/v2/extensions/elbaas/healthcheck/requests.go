package healthcheck

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
	//"github.com/huaweicloud/golangsdk/pagination"
)

// Constants that represent approved monitoring types.
const (
	TypePING  = "PING"
	TypeTCP   = "TCP"
	TypeHTTP  = "HTTP"
	TypeHTTPS = "HTTPS"
)

var (
	errDelayMustGETimeout = fmt.Errorf("Delay must be greater than or equal to timeout")
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToHealthCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Required.  Specifies the ID of the listener to which the health check task belongs.
	ListenerID string `json:"listener_id" required:"true"`
	// Optional. Specifies the protocol used for the health check. The value can be HTTP or TCP (case-insensitive).
	HealthcheckProtocol string `json:"healthcheck_protocol,omitempty"`
	// Optional. Specifies the URI for health check. This parameter is valid when healthcheck_ protocol is HTTP.
	// The value is a string of 1 to 80 characters that must start with a slash (/) and can only contain letters, digits,
	// and special characters, such as -/.%?#&.
	HealthcheckUri string `json:"healthcheck_uri,omitempty"`
	// Optional. Specifies the port used for the health check.  The value ranges from 1 to 65535.
	HealthcheckConnectPort int `json:"healthcheck_connect_port,omitempty"`
	// Optional. MSpecifies the threshold at which the health check result is success, that is, the number of consecutive successful
	// health checks when the health check result of the backend server changes from fail to success.
	// The value ranges from 1 to 10.
	HealthyThreshold int `json:"healthy_threshold,omitempty"`
	// Optional. Specifies the threshold at which the health check result is fail, that is, the number of consecutive
	// failed health checks when the health check result of the backend server changes from success to fail.
	// The value ranges from 1 to 10.
	UnhealthyThreshold int `json:"unhealthy_threshold,omitempty"`
	// Optional. Specifies the maximum timeout duration (s) for the health check.
	// The value ranges from 1 to 50.
	HealthcheckTimeout int `json:"healthcheck_timeout,omitempty"`
	// Optional. Specifies the maximum interval (s) for health check.
	// The value ranges from 1 to 5.
	HealthcheckInterval int `json:"healthcheck_interval,omitempty"`
}

// ToHealthCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToHealthCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
 Create is an operation which provisions a new Health Monitor. There are
 different types of Monitor you can provision: PING, TCP or HTTP(S). Below
 are examples of how to create each one.

 Here is an example config struct to use when creating a PING or TCP Monitor:

 CreateOpts{Type: TypePING, Delay: 20, Timeout: 10, MaxRetries: 3}
 CreateOpts{Type: TypeTCP, Delay: 20, Timeout: 10, MaxRetries: 3}

 Here is an example config struct to use when creating a HTTP(S) Monitor:

 CreateOpts{Type: TypeHTTP, Delay: 20, Timeout: 10, MaxRetries: 3,
 HttpMethod: "HEAD", ExpectedCodes: "200", PoolID: "2c946bfc-1804-43ab-a2ff-58f6a762b505"}
*/
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToHealthCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves a particular Health Monitor based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToHealthUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Optional. Specifies the protocol used for the health check. The value can be HTTP or TCP (case-insensitive).
	HealthcheckProtocol string `json:"healthcheck_protocol,omitempty"`
	// Optional. Specifies the URI for health check. This parameter is valid when healthcheck_ protocol is HTTP.
	// The value is a string of 1 to 80 characters that must start with a slash (/) and can only contain letters, digits,
	// and special characters, such as -/.%?#&.
	HealthcheckUri string `json:"healthcheck_uri,omitempty"`
	// Optional. Specifies the port used for the health check.  The value ranges from 1 to 65535.
	HealthcheckConnectPort int `json:"healthcheck_connect_port,omitempty"`
	// Optional. MSpecifies the threshold at which the health check result is success, that is, the number of consecutive successful
	// health checks when the health check result of the backend server changes from fail to success.
	// The value ranges from 1 to 10.
	HealthyThreshold int `json:"healthy_threshold,omitempty"`
	// Optional. Specifies the threshold at which the health check result is fail, that is, the number of consecutive
	// failed health checks when the health check result of the backend server changes from success to fail.
	// The value ranges from 1 to 10.
	UnhealthyThreshold int `json:"unhealthy_threshold,omitempty"`
	// Optional. Specifies the maximum timeout duration (s) for the health check.
	// The value ranges from 1 to 50.
	HealthcheckTimeout int `json:"healthcheck_timeout,omitempty"`
	// Optional. Specifies the maximum interval (s) for health check.
	// The value ranges from 1 to 5.
	HealthcheckInterval int `json:"healthcheck_interval,omitempty"`
}

// ToHealthpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToHealthUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified Monitor.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToHealthUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will permanently delete a particular Health based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}
