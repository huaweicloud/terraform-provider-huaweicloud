package healthcheck

import (
	"github.com/huaweicloud/golangsdk"
	//"github.com/huaweicloud/golangsdk/pagination"
	//"fmt"
)

// Health represents a load balancer health check. A health monitor is used
// to determine whether or not back-end members of the VIP's pool are usable
// for processing a request. A pool can have several health monitors associated
// with it. There are different types of health monitors supported:
//
// PING: used to ping the members using ICMP.
// TCP: used to connect to the members using TCP.
// HTTP: used to send an HTTP request to the member.
// HTTPS: used to send a secure HTTP request to the member.
//
// When a pool has several monitors associated with it, each member of the pool
// is monitored by all these monitors. If any monitor declares the member as
// unhealthy, then the member status is changed to INACTIVE and the member
// won't participate in its pool's load balancing. In other words, ALL monitors
// must declare the member to be healthy for it to stay ACTIVE.
type Health struct {
	// Specifies the maximum interval (s) for health check.
	HealthcheckInterval int `json:"healthcheck_interval,omitempty"`
	// Specifies the ID of the listener to which the health check task belongs.
	ListenerID string `json:"listener_id" required:"true"`
	// Specifies the health check ID.
	ID string `json:"id"`
	// Specifies the protocol used for the health check. The value can be HTTP or TCP (case-insensitive).
	HealthcheckProtocol string `json:"healthcheck_protocol,omitempty"`
	// Specifies the threshold at which the health check result is fail, that is, the number of consecutive
	// failed health checks when the health check result of the backend server changes from success to fail.
	UnhealthyThreshold int `json:"unhealthy_threshold,omitempty"`
	// Specifies the time when information about the health check was updated.
	UpdateTime string `json:"update_time"`
	// Specifies the time when the health check was created.
	CreateTime string `json:"create_time"`
	// Specifies the port used for the health check.  The value ranges from 1 to 65535.
	HealthcheckConnectPort int `json:"healthcheck_connect_port,omitempty"`
	// Specifies the maximum timeout duration (s) for the health check.
	HealthcheckTimeout int `json:"healthcheck_timeout,omitempty"`
	// Specifies the URI for health check.
	HealthcheckUri string `json:"healthcheck_uri,omitempty"`
	// Specifies the threshold at which the health check result is success, that is, the number of consecutive successful
	// health checks when the health check result of the backend server changes from fail to success.
	HealthyThreshold int `json:"healthy_threshold,omitempty"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a monitor.
func (r commonResult) Extract() (*Health, error) {
	//fmt.Printf("Extracting Health...\n")
	l := new(Health)
	err := r.ExtractInto(l)
	if err != nil {
		//fmt.Printf("Error: %s.\n", err.Error())
		return nil, err
	} else {
		//fmt.Printf("Returning extract: %+v.\n", l)
		return l, nil
	}
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}
