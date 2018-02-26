package listeners

import (
	"github.com/huaweicloud/golangsdk"
)

type CreateResponse struct {
	UpdateTime         string   `json:"update_time"`
	BackendPort        int      `json:"backend_port"`
	ID                 string   `json:"id"`
	BackendProtocol    string   `json:"backend_protocol"`
	StickySessionType  string   `json:"sticky_session_type"`
	Description        string   `json:"description"`
	LoadbalancerID     string   `json:"loadbalancer_id"`
	CreateTime         string   `json:"create_time"`
	Status             string   `json:"status"`
	Protocol           string   `json:"protocol"`
	Port               int      `json:"port"`
	CookieTimeout      int      `json:"cookie_timeout"`
	AdminStateUp       bool     `json:"admin_state_up"`
	SessionSticky      bool     `json:"session_sticky"`
	LbAlgorithm        string   `json:"lb_algorithm"`
	Name               string   `json:"name"`
	TcpDraining        bool     `json:"tcp_draining"`
	TcpDrainingTimeout int      `json:"tcp_draining_timeout"`
	SslProtocols       string   `json:"ssl_protocols"`
	SslCiphers         string   `json:"ssl_ciphers"`
	CertificateID      string   `json:"certificate_id"`
	Certificates       []string `json:"certificates"`
}

type UpdateResponse struct {
	UpdateTime         string   `json:"update_time"`
	BackendPort        int      `json:"backend_port"`
	ID                 string   `json:"id"`
	BackendProtocol    string   `json:"backend_protocol"`
	StickySessionType  string   `json:"sticky_session_type"`
	Description        string   `json:"description"`
	LoadbalancerID     string   `json:"loadbalancer_id"`
	CreateTime         string   `json:"create_time"`
	Status             string   `json:"status"`
	Protocol           string   `json:"protocol"`
	Port               int      `json:"port"`
	CookieTimeout      int      `json:"cookie_timeout"`
	AdminStateUp       bool     `json:"admin_state_up"`
	HealthcheckID      string   `json:"healthcheck_id"`
	SessionSticky      bool     `json:"session_sticky"`
	LbAlgorithm        string   `json:"lb_algorithm"`
	Name               string   `json:"name"`
	TcpDraining        bool     `json:"tcp_draining"`
	TcpDrainingTimeout int      `json:"tcp_draining_timeout"`
	CertificateID      string   `json:"certificate_id"`
	Certificates       []string `json:"certificates"`
}

type Listener struct {
	UpdateTime         string   `json:"update_time"`
	BackendPort        int      `json:"backend_port"`
	ID                 string   `json:"id"`
	BackendProtocol    string   `json:"backend_protocol"`
	StickySessionType  string   `json:"sticky_session_type"`
	Description        string   `json:"description"`
	LoadbalancerID     string   `json:"loadbalancer_id"`
	CreateTime         string   `json:"create_time"`
	Status             string   `json:"status"`
	Protocol           string   `json:"protocol"`
	Port               int      `json:"port"`
	CookieTimeout      int      `json:"cookie_timeout"`
	AdminStateUp       bool     `json:"admin_state_up"`
	MemberNumber       int      `json:"member_number"`
	HealthcheckID      string   `json:"healthcheck_id"`
	SessionSticky      bool     `json:"session_sticky"`
	LbAlgorithm        string   `json:"lb_algorithm"`
	Name               string   `json:"name"`
	CertificateID      string   `json:"certificate_id"`
	Certificates       []string `json:"certificates"`
	TcpDraining        bool     `json:"tcp_draining"`
	TcpDrainingTimeout int      `json:"tcp_draining_timeout"`
	TcpTimeout         int      `json:"tcp_timeout"`
	UdpTimeout         int      `json:"udp_timeout"`
	SslProtocols       string   `json:"ssl_protocols"`
	SslCiphers         string   `json:"ssl_ciphers"`
}

type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (*CreateResponse, error) {
	s := &CreateResponse{}
	return s, r.ExtractInto(s)
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*Listener, error) {
	s := &Listener{}
	return s, r.ExtractInto(s)
}

type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) Extract() (*UpdateResponse, error) {
	s := &UpdateResponse{}
	return s, r.ExtractInto(s)
}
