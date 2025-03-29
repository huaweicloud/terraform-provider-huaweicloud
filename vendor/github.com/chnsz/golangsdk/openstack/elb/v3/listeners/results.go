package listeners

import (
	"github.com/chnsz/golangsdk"
)

type LoadBalancerID struct {
	ID string `json:"id"`
}

// Listener is the primary load balancing configuration object that specifies
// the loadbalancer and port on which client traffic is received, as well
// as other details such as the load balancing method to be use, protocol, etc.
type Listener struct {
	// The unique ID for the Listener.
	ID string `json:"id"`

	// The administrative state of the Listener. A valid value is true (UP) or false (DOWN).
	AdminStateUp bool `json:"admin_state_up"`

	// the ID of the CA certificate used by the listener.
	CAContainerRef string `json:"client_ca_tls_container_ref"`

	// The maximum number of connections allowed for the Loadbalancer.
	// Default is -1, meaning no limit.
	ConnLimit int `json:"connection_limit"`

	// The UUID of default pool. Must have compatible protocol with listener.
	DefaultPoolID string `json:"default_pool_id"`

	// A reference to a Barbican container of TLS secrets.
	DefaultTlsContainerRef string `json:"default_tls_container_ref"`

	// Human-readable description for the Listener.
	Description string `json:"description"`

	// whether to use HTTP2.
	Http2Enable bool `json:"http2_enable"`

	// A list of load balancer IDs.
	Loadbalancers []LoadBalancerID `json:"loadbalancers"`

	// Human-readable name for the Listener. Does not have to be unique.
	Name string `json:"name"`

	// The protocol to loadbalance. A valid value is TCP, HTTP, or HTTPS.
	Protocol string `json:"protocol"`

	// The port on which to listen to client traffic that is associated with the
	// Loadbalancer. A valid value is from 0 to 65535.
	ProtocolPort int `json:"protocol_port"`

	// The list of references to TLS secrets.
	SniContainerRefs []string `json:"sni_container_refs"`

	// Specifies the security policy used by the listener.
	TlsCiphersPolicy string `json:"tls_ciphers_policy"`

	// Whether enable member retry
	EnableMemberRetry bool `json:"enable_member_retry"`

	// Whether enable proxy protocol
	ProxyProtocolEnable bool `json:"proxy_protocol_enable"`

	// The keepalive timeout of the Listener.
	KeepaliveTimeout int `json:"keepalive_timeout"`

	// The client timeout of the Listener.
	ClientTimeout int `json:"client_timeout"`

	// The maximum number of concurrent connections that a listener can handle per second.
	Connection int `json:"connection"`

	// The maximum number of new connections that a listener can handle per second.
	Cps int `json:"cps"`

	// The member timeout of the Listener.
	MemberTimeout int `json:"member_timeout"`

	// The ipgroup of the Listener.
	IpGroup IpGroup `json:"ipgroup"`

	// The http insert headers of the Listener.
	InsertHeaders InsertHeaders `json:"insert_headers"`

	// Transparent client ip enable
	TransparentClientIP bool `json:"transparent_client_ip_enable"`

	// The UUID of the enterprise project who owns the Loadbalancer.
	EnterpriseProjectID string `json:"enterprise_project_id"`

	// The creation time of the current listener
	CreatedAt string `json:"created_at"`

	// The update time of the current listener
	UpdatedAt string `json:"updated_at"`

	// The port range of the current listener
	PortRanges []PortRange `json:"port_ranges"`

	// ELB gzip enable
	GzipEnable bool `json:"gzip_enable"`

	// The QUIC configuration for the current listener
	QuicConfig QuicConfig `json:"quic_config"`

	// Security Policy ID
	SecurityPolicyId string `json:"security_policy_id"`

	// The SNI certificates used by the listener.
	SniMatchAlgo string `json:"sni_match_algo"`

	// Whether enable ssl early data
	SslEarlyDataEnable bool `json:"ssl_early_data_enable"`

	// Enhance L7policy enable
	EnhanceL7policy bool `json:"enhance_l7policy_enable"`

	// Update protection status
	ProtectionStatus string `json:"protection_status"`

	// Update protection reason
	ProtectionReason string `json:"protection_reason"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a listener.
func (r commonResult) Extract() (*Listener, error) {
	var s struct {
		Listener *Listener `json:"listener"`
	}
	err := r.ExtractInto(&s)
	return s.Listener, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Listener.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Listener.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Listener.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
