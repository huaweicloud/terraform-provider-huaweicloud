package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/auth"
)

// CreateOpts is the structure required by the Create method to create a new dedicated microservice instance.
type CreateOpts struct {
	// Host information.
	HostName string `json:"hostName" required:"true"`
	// Access address information.
	Endpoints []string `json:"endpoints" required:"true"`
	// Microservice version.
	Version string `json:"version,omitempty"`
	// Instance status. Value: UP, DOWN, STARTING, or OUTOFSERVICE. Default value: UP.
	Status string `json:"status,omitempty"`
	// Extended attribute. You can customize a key and value. The value must be at least 1 byte long.
	Properties map[string]interface{} `json:"properties,omitempty"`
	// Health check information.
	HealthCheck *HealthCheck `json:"healthCheck,omitempty"`
	// Data center information.
	DataCenterInfo *DataCenter `json:"dataCenterInfo,omitempty"`
}

// HealthCheck is an object that specifies the configuration of the instance health check.
type HealthCheck struct {
	// Heartbeat mode. Value: push or pull.
	Mode string `json:"mode" required:"true"`
	// Heartbeat interval. Unit: s. If the value is less than 5s, the registration is performed at an interval of 5s.
	Interval int `json:"interval" required:"true"`
	// Maximum retries.
	Times int `json:"times" required:"true"`
	// Port.
	Port int `json:"port,omitempty"`
}

// DataCenter is an object that specifies the configuration of the instance data center.
type DataCenter struct {
	// Region name.
	Name string `json:"name" required:"true"`
	// Region.
	Region string `json:"region" required:"true"`
	// AZ.
	AvailableZone string `json:"availableZone" required:"true"`
}

// Create is a method to create a dedicated microservice instance using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts, serviceId, token string) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "instance")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(rootURL(c, serviceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: auth.BuildMoreHeaderUsingToken(c, token),
	})
	return &r, err
}

// Get is a method to retrieves a particular configuration based on its unique ID (and token).
func Get(c *golangsdk.ServiceClient, serviceId, instanceId, token string) (*Instance, error) {
	var r struct {
		Instance Instance `json:"instance"`
	}
	_, err := c.Get(resourceURL(c, serviceId, instanceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: auth.BuildMoreHeaderUsingToken(c, token),
	})
	return &r.Instance, err
}

// List is a method to retrieves all instance informations under specified microservice using its related microservice
// ID (and token).
func List(c *golangsdk.ServiceClient, serviceId string, token string) ([]Instance, error) {
	var r struct {
		Instances []Instance `json:"instances"`
	}
	_, err := c.Get(rootURL(c, serviceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: auth.BuildMoreHeaderUsingToken(c, token),
	})
	return r.Instances, err
}

// Delete is a method to remove an existing instance using its related microservice ID, instance ID (and token).
func Delete(c *golangsdk.ServiceClient, serviceId, instanceId, token string) error {
	_, err := c.Delete(resourceURL(c, serviceId, instanceId), &golangsdk.RequestOpts{
		MoreHeaders: auth.BuildMoreHeaderUsingToken(c, token),
	})
	return err
}
