package policies

import (
	"github.com/chnsz/golangsdk"
)

// Policy contains the infomateion of the policy.
type Policy struct {
	Id            string       `json:"id"`
	Name          string       `json:"name"`
	Action        Action       `json:"action"`
	Options       PolicyOption `json:"options"`
	Level         int          `json:"level"`
	FullDetection bool         `json:"full_detection"`
	BindHosts     []BindHost   `json:"bind_host"`
}

//Action contains actions after the attack is detected
type Action struct {
	Category         string `json:"category,omitempty"`
	FollowedActionId string `json:"followed_action_id,omitempty"`
}

// PolicyOption contains the protection rule of a policy
type PolicyOption struct {
	Webattack      bool `json:"webattack,omitempty"`
	Common         bool `json:"common,omitempty"`
	Crawler        bool `json:"crawler,omitempty"`
	CrawlerEngine  bool `json:"crawler_engine,omitempty"`
	CrawlerScanner bool `json:"crawler_scanner,omitempty"`
	CrawlerScript  bool `json:"crawler_script,omitempty"`
	CrawlerOther   bool `json:"crawler_other,omitempty"`
	Webshell       bool `json:"webshell,omitempty"`
	Cc             bool `json:"cc,omitempty"`
	Custom         bool `json:"custom,omitempty"`
	Whiteblackip   bool `json:"whiteblackip,omitempty"`
	Ignore         bool `json:"ignore,omitempty"`
	Privacy        bool `json:"privacy,omitempty"`
	Antitamper     bool `json:"antitamper,omitempty"`
}

// BindHost the hosts bound to this policy.
type BindHost struct {
	Id       string `json:"id,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	WafType  string `json:"waf_type,omitempty"`
	Mode     string `json:"mode,omitempty"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a policy.
func (r commonResult) Extract() (*Policy, error) {
	var response Policy
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Policy.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Policy.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Policy.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// ListPolicyRst
type ListPolicyRst struct {
	// total policy count.
	Total int `json:"total"`
	// the policy list
	Items []Policy `json:"items"`
}
