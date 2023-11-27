package connections

import "github.com/chnsz/golangsdk/pagination"

type Connections struct {
	ID                  string       `json:"id"`
	Name                string       `json:"name"`
	Status              string       `json:"status"`
	VgwId               string       `json:"vgw_id"`
	VgwIp               string       `json:"vgw_ip"`
	Style               string       `json:"style"`
	CgwId               string       `json:"cgw_id"`
	PeerSubnets         []string     `json:"peer_subnets"`
	TunnelLocalAddress  string       `json:"tunnel_local_address"`
	TunnelPeerAddress   string       `json:"tunnel_peer_address"`
	EnableNqa           bool         `json:"enable_nqa"`
	PolicyRules         []PolicyRule `json:"policy_rules"`
	Ikepolicy           IkePolicy    `json:"ikepolicy"`
	Ipsecpolicy         IpsecPolicy  `json:"ipsecpolicy"`
	CreatedAt           string       `json:"created_at"`
	UpdatedAt           string       `json:"updated_at"`
	EnterpriseProjectID string       `json:"enterprise_project_id"`
	ConnectionMonitorID string       `json:"connection_monitor_id"`
	HaRole              string       `json:"ha_role"`
}

type listResp struct {
	Connections []Connections `json:"vpn_connections"`
	RequestId   string        `json:"request_id"`
	PageInfo    pageInfo      `json:"page_info"`
	TotalCount  int64         `json:"total_count"`
}

type PolicyRule struct {
	RuleIndex   int      `json:"rule_index"`
	Source      string   `json:"source"`
	Destination []string `json:"destination"`
}

type IkePolicy struct {
	IkeVersion              string `json:"ike_version"`
	Phase1NegotiationMode   string `json:"phase1_negotiation_mode"`
	AuthenticationAlgorithm string `json:"authentication_algorithm"`
	EncryptionAlgorithm     string `json:"encryption_algorithm"`
	DhGroup                 string `json:"dh_group"`
	AuthenticationMethod    string `json:"authentication_method"`
	LifetimeSeconds         int    `json:"lifetime_seconds"`
	LocalIdType             string `json:"local_id_type"`
	LocalId                 string `json:"local_id"`
	PeerIdType              string `json:"peer_id_type"`
	PeerId                  string `json:"peer_id"`
	Dpd                     Dpd    `json:"dpd"`
}

type Dpd struct {
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
	Msg      string `json:"msg"`
}

type IpsecPolicy struct {
	AuthenticationAlgorithm string `json:"authentication_algorithm"`
	EncryptionAlgorithm     string `json:"encryption_algorithm"`
	Pfs                     string `json:"pfs"`
	TransformProtocol       string `json:"transform_protocol"`
	LifetimeSeconds         int64  `json:"lifetime_seconds"`
	EncapsulationMode       string `json:"encapsulation_mode"`
}

type pageInfo struct {
	NextMarker   string `json:"next_marker"`
	CurrentCount int    `json:"current_count"`
}

type ConnectionsPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult is empty.
func (r ConnectionsPage) IsEmpty() (bool, error) {
	resp, err := extractConnections(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r ConnectionsPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	return resp.NextMarker, nil
}

// extractConnections is a method which to extract the response to a VpnConnection list.
func extractConnections(r pagination.Page) ([]Connections, error) {
	var s listResp
	err := r.(ConnectionsPage).Result.ExtractInto(&s)
	return s.Connections, err
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listResp
	err := r.(ConnectionsPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}
