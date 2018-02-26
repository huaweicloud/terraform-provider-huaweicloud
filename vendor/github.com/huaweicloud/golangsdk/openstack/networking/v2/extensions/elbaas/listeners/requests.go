package listeners

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Protocol string

// Supported attributes for create/update operations.
const (
	ProtocolTCP   Protocol = "TCP"
	ProtocolHTTP  Protocol = "HTTP"
	ProtocolHTTPS Protocol = "HTTPS"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToListenerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular listener attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	LoadbalancerId string `q:"loadbalancer_id"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToListenerListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// routers. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those routers that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToListenerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ListenerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToListenerCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Required.  Specifies the load balancer name.
	// The name is a string of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name" required:"true"`
	// Optional. Provides supplementary information about the listener.
	// The value is a string of 0 to 128 characters and cannot contain angle brackets (<>).
	Description string `json:"description,omitempty"`
	// Required.  Specifies the ID of the load balancer to which the listener belongs.
	LoadbalancerID string `json:"loadbalancer_id" required:"true"`
	// Required.  Specifies the listening protocol used for layer 4 or 7.
	// A listener using UDP is not allowed for a private network load balancer.
	// The value can be HTTP, TCP, HTTPS, SSL, or UDP.
	// The protocol - can either be TCP, HTTP or HTTPS.
	Protocol Protocol `json:"protocol" required:"true"`
	// Required.  Specifies the listening port.
	// The value ranges from 1 to 65535.
	ProtocolPort int `json:"port" required:"true"`
	// Required.  Specifies the backend protocol.
	// If the value of protocol is UDP, the parameter value can only be UDP.
	// If the value of protocol is SSL, the parameter value can only be TCP.
	BackendProtocol Protocol `json:"backend_protocol" required:"true"`
	// Required.  Specifies the backend port.
	// The value ranges from 1 to 65535.
	BackendProtocolPort int `json:"backend_port" required:"true"`
	// Required.  Specifies the load balancing algorithm for the listener.
	// The value can be roundrobin, leastconn, or source.
	Algorithm string `json:"lb_algorithm" required:"true"`
	// Optional.  Specifies whether to enable sticky session.
	// The value can be true or false. Sticky session is enabled when the value is true.
	// If the value of protocol is SSL, the sticky session is not supported and the parameter is invalid.
	// If the value of protocol is HTTP, HTTPS, or TCP and the value of lb_algorithm is not roundrobin,
	// the parameter value can only be false.
	SessionSticky bool `json:"session_sticky,omitempty"`
	// Optional.  Specifies the cookie processing method. The value is insert.
	// insert indicates that the cookie is inserted by the load balancer.
	// This parameter is valid when protocol is set to HTTP and session_sticky to true.
	// The default value is insert. This parameter is invalid when protocol is set to TCP, SSL, or UDP,
	// which means the parameter is empty.
	// The ID of the default pool with which the Listener is associated.
	StickySessionType string `json:"sticky_session_type,omitempty"`
	// Optional.  Specifies the cookie timeout period (minutes). This parameter is valid when protocol is set to HTTP,
	// session_sticky to true, and sticky_session_type to insert. This parameter is invalid when protocol is set to
	// TCP, UDP, or SSL. The value ranges from 1 to 1440.
	CookieTimeout int `json:"cookie_timeout,omitempty"`
	// Optional.  Specifies the TCP timeout period (minutes). This parameter is valid when protocol is set to TCP.
	// The value ranges from 1 to 5.
	TcpTimeout int `json:"tcp_timeout,omitempty"`
	// Optional.  Specifies whether to maintain the TCP connection to the backend ECS after the ECS is deleted.
	// This parameter is valid when protocol is set to TCP.
	// The value can be true or false.
	TcpDraining bool `json:"tcp_draining,omitempty"`
	// Optional.  Specifies the timeout duration (minutes) for the TCP connection to the backend ECS after the ECS
	// is deleted. This parameter is valid when protocol is set to TCP and tcp_draining to true.
	// The value ranges from 0 to 60.
	TcpDrainingTimeout int `json:"tcp_draining_timeout,omitempty"`
	// Optional.  Specifies the certificate ID. This parameter is mandatory when protocol is set to HTTPS or SSL.
	// The value can be obtained by viewing details of the SSL certificate.
	CertificateID string `json:"certificate_id,omitempty"`
	// Optional.  Specifies the SSL certificate ID list if the value of protocol is HTTPS.
	// This parameter is mandatory in SNI scenarios.
	// This parameter is valid only when the load balancer is a public network load balancer.
	Certificates []string `json:"certificates,omitempty"`
	// Optional.  Specifies the UDP session timeout duration (minutes). This parameter is valid when protocol is set to UDP.
	// The value ranges from 1 to 1440.
	UDPTimeout int `json:"udp_timeout,omitempty"`
	// Optional.  Specifies the SSL protocol standard supported by a listener.
	// This parameter is used for enabling specified encryption protocols and valid only when the value of protocol
	// is set to HTTPS or SSL.  The value is TLSv1.2 or TLSv1.2 TLSv1.1 TLSv1. The default value is TLSv1.2.
	SSLProtocols string `json:"ssl_protocols,omitempty"`
	// Optional.  Specifies the cipher suite of an encryption protocol. This parameter is valid only when the value of protocol is set
	// to HTTPS or SSL. The value is Default, Extended, or Strict.
	// The value of Default is ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-SHA384:
	// ECDHE-RSA-AES128-SHA256.
	// The value of Extended is ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:AES128-SHA256:AES256-SHA256:
	// ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES128-SHA:DHE-RSA-AES128-SHA:
	// ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:AES128-SHA:AES256-SHA:DHE-DSS-AES128-SHA:CAMELLIA128-SHA:
	// EDH-RSA-DES-CBC3-SHA:DES-CBC3-SHA:ECDHE-RSA-RC4-SHA:RC4-SHA:DHE-RSA-AES256-SHA:DHE-DSS-AES256-SHA:
	// DHE-RSA-CAMELLIA256-SHA:DHE-DSS-CAMELLIA256-SHA:CAMELLIA256-SHA:EDH-DSS-DES-CBC3-SHA:DHE-RSA-CAMELLIA128-SHA:
	// DHE-DSS-CAMELLIA128-SHA.
	// The value of Strict is ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256.
	// The default value is Default. The value can only be set to Extended if the value of ssl_protocols is set to
	// TLSv1.2 TLSv1.1 TLSv1.
	SSLCiphers string `json:"ssl_ciphers,omitempty"`
}

// ToListenerCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToListenerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is an operation which provisions a new Listeners based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create Listeners on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToListenerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves a particular Listeners based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToListenerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Required.  Specifies the load balancer name.
	// The name is a string of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name,omitempty"`
	// Optional. Provides supplementary information about the listener.
	// The value is a string of 0 to 128 characters and cannot contain angle brackets (<>).
	Description string `json:"description,omitempty"`
	// Required.  Specifies the listening port.
	// The value ranges from 1 to 65535.
	ProtocolPort int `json:"port,omitempty"`
	// Required.  Specifies the backend port.
	// The value ranges from 1 to 65535.
	BackendProtocolPort int `json:"backend_port,omitempty"`
	// Required.  Specifies the load balancing algorithm for the listener.
	// The value can be roundrobin, leastconn, or source.
	Algorithm string `json:"lb_algorithm,omitempty"`
	// Optional.  Specifies the TCP timeout period (minutes). This parameter is valid when protocol is set to TCP.
	// The value ranges from 1 to 5.
	TcpTimeout int `json:"tcp_timeout,omitempty"`
	// Optional.  Specifies whether to maintain the TCP connection to the backend ECS after the ECS is deleted.
	// This parameter is valid when protocol is set to TCP.
	// The value can be true or false.
	TcpDraining bool `json:"tcp_draining,omitempty"`
	// Optional.  Specifies the timeout duration (minutes) for the TCP connection to the backend ECS after the ECS
	// is deleted. This parameter is valid when protocol is set to TCP and tcp_draining to true.
	// The value ranges from 0 to 60.
	TcpDrainingTimeout int `json:"tcp_draining_timeout,omitempty"`
	// Optional.  Specifies the UDP session timeout duration (minutes). This parameter is valid when protocol is set to UDP.
	// The value ranges from 1 to 1440.
	UDPTimeout int `json:"udp_timeout,omitempty"`
	// Optional.  Specifies the SSL protocol standard supported by a listener.
	// This parameter is used for enabling specified encryption protocols and valid only when the value of protocol
	// is set to HTTPS or SSL.  The value is TLSv1.2 or TLSv1.2 TLSv1.1 TLSv1. The default value is TLSv1.2.
	SSLProtocols string `json:"ssl_protocols,omitempty"`
	// Optional.  Specifies the cipher suite of an encryption protocol. This parameter is valid only when the value of protocol is set
	// to HTTPS or SSL. The value is Default, Extended, or Strict.
	// The value of Default is ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-SHA384:
	// ECDHE-RSA-AES128-SHA256.
	// The value of Extended is ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:AES128-SHA256:AES256-SHA256:
	// ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES128-SHA:DHE-RSA-AES128-SHA:
	// ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:AES128-SHA:AES256-SHA:DHE-DSS-AES128-SHA:CAMELLIA128-SHA:
	// EDH-RSA-DES-CBC3-SHA:DES-CBC3-SHA:ECDHE-RSA-RC4-SHA:RC4-SHA:DHE-RSA-AES256-SHA:DHE-DSS-AES256-SHA:
	// DHE-RSA-CAMELLIA256-SHA:DHE-DSS-CAMELLIA256-SHA:CAMELLIA256-SHA:EDH-DSS-DES-CBC3-SHA:DHE-RSA-CAMELLIA128-SHA:
	// DHE-DSS-CAMELLIA128-SHA.
	// The value of Strict is ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256.
	// The default value is Default. The value can only be set to Extended if the value of ssl_protocols is set to
	// TLSv1.2 TLSv1.1 TLSv1.
	SSLCiphers string `json:"ssl_ciphers,omitempty"`
}

// ToListenerUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToListenerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified Listener.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToListenerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will permanently delete a particular Listeners based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	url := resourceURL(c, id)
	//fmt.Printf("Delete listener url: %s.\n", url)
	_, r.Err = c.Delete(url, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}
