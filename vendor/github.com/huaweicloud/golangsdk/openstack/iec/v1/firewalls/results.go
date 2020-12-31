package firewalls

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

type RespFirewallRulesEntity struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	Action      string `json:"action,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	IPVersion   int    `json:"ip_version,omitempty"`
	DstIPAddr   string `json:"destination_ip_address,omitempty"`
	DstPort     string `json:"destination_port,omitempty"`
	SrcIPAddr   string `json:"source_ip_address,omitempty"`
	SrcPort     string `json:"source_port,omitempty"`
	OperateType string `json:"operate_type,omitempty"`
}

type RespPolicyEntity struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name,omitempty"`
	Description   string                    `json:"description,omitempty"`
	Audited       bool                      `json:"audited,omitempty"`
	Shared        bool                      `json:"shared,omitempty"`
	FirewallRules []RespFirewallRulesEntity `json:"firewall_rules,omitempty"`
}

type Firewall struct {
	ID                 string           `json:"id"`
	DomainID           string           `json:"domain_id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	AdminStateUp       bool             `json:"admin_state_up"`
	Status             string           `json:"status"`
	IngressFWPolicy    RespPolicyEntity `json:"ingress_firewall_policy"`
	EgressFWPolicy     RespPolicyEntity `json:"egress_firewall_policy"`
	IngressFWRuleCount int64            `json:"ingress_firewall_rule_count,omitempty"`
	EgressFWRuleCount  int64            `json:"egress_firewall_rule_count,omitempty"`
	Subnets            []common.Subnet  `json:"subnets,omitempty"`
}

type UpdateFirewallResp struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*Firewall, error) {
	var entity Firewall
	err := r.ExtractIntoStructPtr(&entity, "firewall")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*Firewall, error) {
	var entity Firewall
	err := r.ExtractIntoStructPtr(&entity, "firewall")
	return &entity, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*UpdateFirewallResp, error) {
	var entity UpdateFirewallResp
	err := r.ExtractIntoStructPtr(&entity, "firewall")
	return &entity, err
}

type UpdateRuleResult struct {
	commonResult
}

func (r UpdateRuleResult) Extract() (*Firewall, error) {
	var entity Firewall
	err := r.ExtractIntoStructPtr(&entity, "firewall")
	return &entity, err
}
