package rules

import (
	"github.com/chnsz/golangsdk"
)

// BatchCreateOpts is the structure required by the BatchCreate method to batch create at least one forward rule.
type BatchCreateOpts struct {
	// Forwarding rules created in batches.
	Rules []RuleOpts `json:"rules,omitempty"`
}

// RuleOpts is the object that represents the forwarding rule configuration structure.
type RuleOpts struct {
	// The forward protocol.
	ForwardProtocol string `json:"forward_protocol" required:"true"`
	// The forward port.
	// The valid value is range from 1 to 65535.
	ForwardPort int `json:"forward_port" required:"true"`
	// The source port.
	// The valid value is range from 1 to 65535.
	SourcePort int `json:"source_port" required:"true"`
	// The source IP list, separate the IPs with commas.
	SourceIp string `json:"source_ip" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// BatchCreate is a method to to batch create at least one forward rule using given parameters.
func BatchCreate(c *golangsdk.ServiceClient, instnaceId, ip string, opts BatchCreateOpts) (*BatchResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r BatchResp
	_, err = c.Post(batchCreateURL(c, instnaceId, ip), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// List is a method to query the list of the forward rule using given opts.
func List(c *golangsdk.ServiceClient, instnaceId, ip string) ([]Rule, error) {
	var r ListResp
	_, err := c.Get(listURL(c, instnaceId, ip), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return r.Rules, err
}

// UpdateOpts is the structure required by the Update method to update the configurations of the forward rule.
type UpdateOpts struct {
	// The forward protocol.
	ForwardProtocol string `json:"forward_protocol" required:"true"`
	// The forward port.
	// The valid value is range from 1 to 65535.
	ForwardPort int `json:"forward_port" required:"true"`
	// The source port.
	// The valid value is range from 1 to 65535.
	SourcePort int `json:"source_port" required:"true"`
	// The source IP list, separate the IPs with commas.
	SourceIp string `json:"source_ip" required:"true"`
}

// Update is a method to update the configurations of the forward rule using given parameters.
func Update(c *golangsdk.ServiceClient, instnaceId, ip, ruleId string, opts UpdateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	var r BatchResp
	_, err = c.Put(updateURL(c, instnaceId, ip, ruleId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// BatchDeleteOpts is the structure required by the BatchDelete method to batch delete at least one forward rule.
type BatchDeleteOpts struct {
	RuleIds []string `json:"ids" required:"true"`
}

// BatchDelete is a method to to batch delete at least one forward rule using given parameters.
func BatchDelete(c *golangsdk.ServiceClient, instnaceId, ip string, opts BatchDeleteOpts) (*BatchResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r BatchResp
	_, err = c.Post(batchDeleteURL(c, instnaceId, ip), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}
