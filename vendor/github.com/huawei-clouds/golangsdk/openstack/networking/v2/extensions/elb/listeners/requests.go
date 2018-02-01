package listeners

import (
	"log"

	"github.com/huawei-clouds/golangsdk"
)

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
	Name               string `json:"name" required:"true"`
	Description        string `json:"description,omitempty"`
	LoadbalancerID     string `json:"loadbalancer_id" required:"true"`
	Protocol           string `json:"protocol" required:"true"`
	Port               int    `json:"port" required:"true"`
	BackendProtocol    string `json:"backend_protocol" required:"true"`
	BackendPort        int    `json:"backend_port" required:"true"`
	LbAlgorithm        string `json:"lb_algorithm" required:"true"`
	SessionSticky      bool   `json:"session_sticky,omitempty"`
	StickySessionType  string `json:"sticky_session_type,omitempty"`
	CookieTimeout      int    `json:"cookie_timeout,omitempty"`
	TcpTimeout         int    `json:"tcp_timeout,omitempty"`
	TcpDraining        bool   `json:"tcp_draining,omitempty"`
	TcpDrainingTimeout int    `json:"tcp_draining_timeout,omitempty"`
	CertificateID      string `json:"certificate_id,omitempty"`
	UdpTimeout         int    `json:"udp_timeout,omitempty"`
	SslProtocols       string `json:"ssl_protocols,omitempty"`
	SslCiphers         string `json:"ssl_ciphers,omitempty"`
}

// ToListenerCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToListenerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClientExtension, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToListenerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create url:%q, body=%#v", rootURL(c), b)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *golangsdk.ServiceClientExtension, id string) (r GetResult) {
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
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	Port               int    `json:"port,omitempty"`
	BackendPort        int    `json:"backend_port,omitempty"`
	LbAlgorithm        string `json:"lb_algorithm,omitempty"`
	TcpTimeout         int    `json:"tcp_timeout,omitempty"`
	TcpDraining        bool   `json:"tcp_draining,omitempty"`
	TcpDrainingTimeout int    `json:"tcp_draining_timeout,omitempty"`
	UdpTimeout         int    `json:"udp_timeout,omitempty"`
	SslProtocols       string `json:"ssl_protocols,omitempty"`
	SslCiphers         string `json:"ssl_ciphers,omitempty"`
	CertificateID      string `json:"certificate_id,omitempty"`
}

func (u UpdateOpts) IsNeedUpdate() (bool, error) {
	d, e := u.ToListenerUpdateMap()
	if e == nil {
		return len(d) != 0, nil
	}
	return false, e
}

// ToListenerUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToListenerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified Listener.
func Update(c *golangsdk.ServiceClientExtension, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToListenerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular Listener based on its unique ID.
func Delete(c *golangsdk.ServiceClientExtension, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
