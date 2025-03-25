package listeners

import (
	"github.com/chnsz/golangsdk"
)

// Type Protocol represents a listener protocol.
type Protocol string

// Supported attributes for create/update operations.
const (
	ProtocolTCP   Protocol = "TCP"
	ProtocolUDP   Protocol = "UDP"
	ProtocolHTTP  Protocol = "HTTP"
	ProtocolHTTPS Protocol = "HTTPS"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToListenerCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options for creating a listener.
type CreateOpts struct {
	// The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// the ID of the CA certificate used by the listener.
	CAContainerRef string `json:"client_ca_tls_container_ref,omitempty"`

	// The ID of the default pool with which the Listener is associated.
	DefaultPoolID string `json:"default_pool_id,omitempty"`

	// A reference to a Barbican container of TLS secrets.
	DefaultTlsContainerRef string `json:"default_tls_container_ref,omitempty"`

	// Human-readable description for the Listener.
	Description string `json:"description,omitempty"`

	// whether to use HTTP2.
	Http2Enable *bool `json:"http2_enable,omitempty"`

	// The load balancer on which to provision this listener.
	LoadbalancerID string `json:"loadbalancer_id" required:"true"`

	// Human-readable name for the Listener. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// ProjectID is only required if the caller has an admin role and wants
	// to create a pool for another project.
	ProjectID string `json:"project_id,omitempty"`

	// The protocol - can either be TCP, HTTP or HTTPS.
	Protocol Protocol `json:"protocol" required:"true"`

	// The port on which to listen for client traffic.
	ProtocolPort int `json:"protocol_port,omitempty"`

	// A list of references to TLS secrets.
	SniContainerRefs []string `json:"sni_container_refs,omitempty"`

	// Specifies the security policy used by the listener.
	TlsCiphersPolicy string `json:"tls_ciphers_policy,omitempty"`

	// Whether enable member retry
	EnableMemberRetry *bool `json:"enable_member_retry,omitempty"`

	// Whether enable proxy protocol
	ProxyProtocolEnable *bool `json:"proxy_protocol_enable,omitempty"`

	// Whether enable ssl early data
	SslEarlyDataEnable *bool `json:"ssl_early_data_enable,omitempty"`

	// The keepalive timeout of the Listener.
	KeepaliveTimeout *int `json:"keepalive_timeout,omitempty"`

	// The client timeout of the Listener.
	ClientTimeout *int `json:"client_timeout,omitempty"`

	// The maximum number of concurrent connections that a listener can handle per second.
	Connection *int `json:"connection,omitempty"`

	// The maximum number of new connections that a listener can handle per second.
	Cps *int `json:"cps,omitempty"`

	// The member timeout of the Listener.
	MemberTimeout *int `json:"member_timeout,omitempty"`

	// The ipgroup of the Listener.
	IpGroup *IpGroup `json:"ipgroup,omitempty"`

	// The http insert headers of the Listener.
	InsertHeaders *InsertHeaders `json:"insert_headers,omitempty"`

	// Transparent client ip enable
	TransparentClientIP *bool `json:"transparent_client_ip_enable,omitempty"`

	// Enhance L7policy enable
	EnhanceL7policy *bool `json:"enhance_l7policy_enable,omitempty"`

	// The port range of the current listener
	PortRanges []PortRange `json:"port_ranges,omitempty"`

	// ELB gzip enable
	GzipEnable *bool `json:"gzip_enable,omitempty"`

	// The QUIC configuration for the current listener
	QuicConfig *QuicConfig `json:"quic_config,omitempty"`

	// Security Policy ID
	SecurityPolicyId string `json:"security_policy_id,omitempty"`

	// The SNI certificates used by the listener.
	SniMatchAlgo string `json:"sni_match_algo,omitempty"`

	// Protection status
	ProtectionStatus string `json:"protection_status,omitempty"`

	// Protection reason
	ProtectionReason string `json:"protection_reason,omitempty"`
}

type IpGroup struct {
	IpGroupId string `json:"ipgroup_id" required:"true"`
	Enable    *bool  `json:"enable_ipgroup,omitempty"`
	Type      string `json:"type,omitempty"`
}

type PortRange struct {
	StartPort int `json:"start_port,omitempty"`
	EndPort   int `json:"end_port,omitempty"`
}

type QuicConfig struct {
	QuicListenerId    string `json:"quic_listener_id" required:"true"`
	EnableQuicUpgrade *bool  `json:"enable_quic_upgrade,omitempty"`
}

type InsertHeaders struct {
	ForwardedELBIP            *bool `json:"X-Forwarded-ELB-IP,omitempty"`
	ForwardedPort             *bool `json:"X-Forwarded-Port,omitempty"`
	ForwardedForPort          *bool `json:"X-Forwarded-For-Port,omitempty"`
	ForwardedHost             *bool `json:"X-Forwarded-Host,omitempty"`
	ForwardedProto            *bool `json:"X-Forwarded-Proto,omitempty"`
	RealIP                    *bool `json:"X-Real-IP,omitempty"`
	ForwardedELBID            *bool `json:"X-Forwarded-ELB-ID,omitempty"`
	ForwardedTLSCertificateID *bool `json:"X-Forwarded-TLS-Certificate-ID,omitempty"`
	ForwardedTLSCipher        *bool `json:"X-Forwarded-TLS-Cipher,omitempty"`
	ForwardedTLSProtocol      *bool `json:"X-Forwarded-TLS-Protocol,omitempty"`
}

// ToListenerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToListenerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "listener")
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
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular Listeners based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToListenerUpdateMap() (map[string]interface{}, error)
}

type IpGroupUpdate struct {
	IpGroupId string `json:"ipgroup_id,omitempty"`
	Type      string `json:"type,omitempty"`
	Enable    *bool  `json:"enable_ipgroup,omitempty"`
}

// UpdateOpts represents options for updating a Listener.
type UpdateOpts struct {
	// Human-readable name for the Listener. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable description for the Listener.
	Description *string `json:"description,omitempty"`

	// The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// the ID of the CA certificate used by the listener.
	CAContainerRef *string `json:"client_ca_tls_container_ref,omitempty"`

	// The ID of the default pool with which the Listener is associated.
	DefaultPoolID string `json:"default_pool_id,omitempty"`

	// A reference to a container of TLS secrets.
	DefaultTlsContainerRef *string `json:"default_tls_container_ref,omitempty"`

	// whether to use HTTP2.
	Http2Enable *bool `json:"http2_enable,omitempty"`

	// A list of references to TLS secrets.
	SniContainerRefs *[]string `json:"sni_container_refs,omitempty"`

	// Specifies the security policy used by the listener.
	TlsCiphersPolicy *string `json:"tls_ciphers_policy,omitempty"`

	// Whether enable member retry
	EnableMemberRetry *bool `json:"enable_member_retry,omitempty"`

	// Whether enable proxy protocol
	ProxyProtocolEnable *bool `json:"proxy_protocol_enable,omitempty"`

	// Whether enable ssl early data
	SslEarlyDataEnable *bool `json:"ssl_early_data_enable,omitempty"`

	// The keepalive timeout of the Listener.
	KeepaliveTimeout *int `json:"keepalive_timeout,omitempty"`

	// The client timeout of the Listener.
	ClientTimeout *int `json:"client_timeout,omitempty"`

	// The maximum number of concurrent connections that a listener can handle per second.
	Connection *int `json:"connection,omitempty"`

	// The maximum number of new connections that a listener can handle per second.
	Cps *int `json:"cps,omitempty"`

	// The member timeout of the Listener.
	MemberTimeout *int `json:"member_timeout,omitempty"`

	// The ipgroup of the Listener.
	IpGroup *IpGroupUpdate `json:"ipgroup,omitempty"`

	// The http insert headers of the Listener.
	InsertHeaders *InsertHeaders `json:"insert_headers,omitempty"`

	// Transparent client ip enable
	TransparentClientIP *bool `json:"transparent_client_ip_enable,omitempty"`

	// Enhance L7policy enable
	EnhanceL7policy *bool `json:"enhance_l7policy_enable,omitempty"`

	// ELB gzip enable
	GzipEnable *bool `json:"gzip_enable,omitempty"`

	// The QUIC configuration for the current listener
	QuicConfig *QuicConfig `json:"quic_config"`

	// Security Policy ID
	SecurityPolicyId *string `json:"security_policy_id"`

	// The SNI certificates used by the listener.
	SniMatchAlgo string `json:"sni_match_algo,omitempty"`

	// Update protection status
	ProtectionStatus string `json:"protection_status,omitempty"`

	// Update protection reason
	ProtectionReason *string `json:"protection_reason,omitempty"`
}

// ToListenerUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToListenerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "listener")
}

// Update is an operation which modifies the attributes of the specified
// Listener.
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
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

// ForceDelete will delete the listener and the sub resource(listener and l7 policies, unbind associated pools)
func ForceDelete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceForceDeleteURL(c, id), nil)
	return
}
