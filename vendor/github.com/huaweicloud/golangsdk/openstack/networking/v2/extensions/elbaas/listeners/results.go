package listeners

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
	//"fmt"
)

type LoadBalancerID struct {
	ID string `json:"id"`
}

// Listener is the primary load balancing configuration object that specifies
// the loadbalancer and port on which client traffic is received, as well
// as other details such as the load balancing method to be use, protocol, etc.
type Listener struct {
	// Specifies the time when information about the listener was updated.
	UpdateTime string `json:"update_time"`
	// Specifies the backend port.
	BackendProtocolPort int `json:"backend_port"`
	// Specifies the listener ID.
	ID string `json:"id"`
	// Specifies the backend protocol.
	BackendProtocol Protocol `json:"backend_protocol"`
	// Specifies the cookie processing method. The value is insert.
	StickySessionType string `json:"sticky_session_type"`
	// Provides supplementary information about the listener.
	Description string `json:"description"`
	// Specifies the ID of the load balancer to which the listener belongs.
	LoadbalancerID string `json:"loadbalancer_id"`
	// Specifies the time when the listener was created.
	CreateTime string `json:"create_time"`
	// Specifies the listener status. The value can be ACTIVE, PENDING_CREATE, or ERROR.
	Status string `json:"status"`
	// Specifies the listening protocol used for layer 4 or 7.
	Protocol Protocol `json:"protocol"`
	// Specifies the listening port.
	ProtocolPort int `json:"port"`
	// Specifies the cookie timeout period (minutes).
	CookieTimeout int `json:"cookie_timeout"`
	// Specifies the status of the load balancer.
	AdminStateUp bool `json:"admin_state_up"`
	// Specifies the quantity of backend ECSs.
	MemberNumber int `json:"member_number"`
	// Specifies the health check task ID.
	HealthCheckID string `json:"healthcheck_id"`
	// Specifies whether to enable sticky session.
	SessionSticky bool `json:"session_sticky,omitempty"`
	// Specifies the load balancing algorithm for the listener.
	Algorithm string `json:"lb_algorithm"`
	// Specifies the load balancer name.
	Name string `json:"name"`
	// Specifies the certificate ID.
	CertificateID string `json:"certificate_id"`
	// Specifies the SSL certificate ID list if the value of protocol is HTTPS.
	Certificates []string `json:"certificates"`
	// Specifies the TCP timeout period (minutes).
	TcpTimeout int `json:"tcp_timeout"`
	// Specifies the UDP session timeout duration (minutes).
	UDPTimeout int `json:"udp_timeout"`
	// Specifies the SSL protocol standard supported by a listener.
	SSLProtocols string `json:"ssl_protocols"`
	// Specifies the cipher suite of an encryption protocol.
	SSLCiphers string `json:"ssl_ciphers"`
	// Secifies whether to maintain the TCP connection to the backend ECS after the ECS is deleted.
	TcpDraining bool `json:"tcp_draining"`
	// Specifies the timeout duration (minutes) for the TCP connection to the backend ECS after the ECS
	TcpDrainingTimeout int `json:"tcp_draining_timeout"`
}

// ListenerPage is the page returned by a pager when traversing over a
// collection of routers.
type ListenerPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ListenerPage) NextPageURL() (string, error) {
	return "", nil
}

// IsEmpty checks whether a RouterPage struct is empty.
func (r ListenerPage) IsEmpty() (bool, error) {
	is, err := ExtractListeners(r)
	return len(is) == 0, err
}

// ExtractListeners accepts a Page struct, specifically a ListenerPage struct,
// and extracts the elements into a slice of Listener structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractListeners(r pagination.Page) ([]Listener, error) {
	var Listeners []Listener
	err := (r.(ListenerPage)).ExtractInto(&Listeners)
	return Listeners, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a router.
func (r commonResult) Extract() (*Listener, error) {
	//fmt.Printf("Extracting Listener...\n")
	l := new(Listener)
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
