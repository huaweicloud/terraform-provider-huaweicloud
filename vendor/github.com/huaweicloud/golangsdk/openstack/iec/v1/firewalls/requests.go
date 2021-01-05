package firewalls

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
)

type CreateOpts struct {
	Name        string `json:"name,true"`
	Description string `json:"description,omitempty"`
}

type CreateOptsBuilder interface {
	ToFirewallCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToFirewallCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "firewall")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFirewallCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(CreateURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, firewallID string) (r DeleteResult) {
	url := DeleteURL(client, firewallID)
	_, r.Err = client.Delete(url, nil)
	return
}

func Get(client *golangsdk.ServiceClient, fwID string) (r GetResult) {
	url := GetURL(client, fwID)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

type ReqSubnet struct {
	ID    string `json:"id"`
	VpcID string `json:"vpc_id"`
}

type UpdateOpts struct {
	Name         string       `json:"name,omitempty"`
	Description  *string      `json:"description,omitempty"`
	AdminStateUp *bool        `json:"admin_state_up,omitempty"`
	Subnets      *[]ReqSubnet `json:"subnets,omitempty"`
}

type UpdateOptsBuilder interface {
	ToUpdateFirewallMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToUpdateFirewallMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "firewall")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, fwID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateFirewallMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UpdateURL(client, fwID), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}

type ReqFirewallRulesOpts struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Name        string `json:"name"`
	Protocol    string `json:"protocol"`
	Action      string `json:"action"`
	IPVersion   int    `json:"ip_version,omitempty"`
	DstIPAddr   string `json:"destination_ip_address"`
	DstPort     string `json:"destination_port"`
	SrcIPAddr   string `json:"source_ip_address"`
	SrcPort     string `json:"source_port"`
	OperateType string `json:"operate_type"`
}

type ReqPolicyOpts struct {
	PolicyID      string                  `json:"id"`
	FirewallRules *[]ReqFirewallRulesOpts `json:"firewall_rules,omitempty"`
}

type UpdateRuleOpts struct {
	ReqFirewallOutPolicy *ReqPolicyOpts `json:"egress_firewall_policy,omitempty"`
	ReqFirewallInPolicy  *ReqPolicyOpts `json:"ingress_firewall_policy,omitempty"`
}

type UpdateRuleOptsBuilder interface {
	ToUpdateFirewallRuleMap() (map[string]interface{}, error)
}

func (opts UpdateRuleOpts) ToUpdateFirewallRuleMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "firewall")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateRule(client *golangsdk.ServiceClient, fwID string, opts UpdateRuleOptsBuilder) (r UpdateRuleResult) {
	b, err := opts.ToUpdateFirewallRuleMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UpdateRuleURL(client, fwID), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}
