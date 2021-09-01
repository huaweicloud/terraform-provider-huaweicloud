package acl

import (
	"github.com/chnsz/golangsdk"
)

// ConsoleACLPolicyBuilder allows extensions to add additional parameters to
// the modify request.
type ConsoleACLPolicyBuilder interface {
	ToConsoleACLPolicyMap() (map[string]interface{}, error)
}

// APIACLPolicyBuilder allows extensions to add additional parameters to
// the modify request.
type APIACLPolicyBuilder interface {
	ToAPIACLPolicyMap() (map[string]interface{}, error)
}

// ACLPolicy provides options used to create, update or get a identity acl.
type ACLPolicy struct {
	AllowAddressNetmasks []AllowAddressNetmasks `json:"allow_address_netmasks,omitempty"`
	AllowIPRanges        []AllowIPRanges        `json:"allow_ip_ranges,omitempty"`
}

// AllowAddressNetmasks provides options for creating, updating or getting a IPv4 CIDR blocks.
type AllowAddressNetmasks struct {
	AddressNetmask string `json:"address_netmask" required:"true"`
	Description    string `json:"description,omitempty"`
}

// AllowIPRanges provides options for creating, updating or getting a IP address ranges.
type AllowIPRanges struct {
	IPRange     string `json:"ip_range" required:"true"`
	Description string `json:"description,omitempty"`
}

// ToConsoleACLPolicyMap formats a create or update opts into a modify request for console access.
func (aclPolicy ACLPolicy) ToConsoleACLPolicyMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(aclPolicy, "console_acl_policy")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ToAPIACLPolicyMap formats a create or update opts into a modify request for api access.
func (aclPolicy ACLPolicy) ToAPIACLPolicyMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(aclPolicy, "api_acl_policy")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ConsoleACLPolicyUpdate can creates a new acl or updates a exist acl for console access.
func ConsoleACLPolicyUpdate(client *golangsdk.ServiceClient, opts ConsoleACLPolicyBuilder, domainID string) (r ACLResult) {
	b, err := opts.ToConsoleACLPolicyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(consoleACLPolicyURL(client, domainID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// APIACLPolicyUpdate can creates a new acl or updates a exist acl for api access.
func APIACLPolicyUpdate(client *golangsdk.ServiceClient, opts APIACLPolicyBuilder, domainID string) (r ACLResult) {
	b, err := opts.ToAPIACLPolicyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(apiACLPolicyURL(client, domainID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// ConsoleACLPolicyGet retrieves details on iam identity acl for console access, by domain ID.
func ConsoleACLPolicyGet(client *golangsdk.ServiceClient, domainID string) (r ACLResult) {
	_, r.Err = client.Get(consoleACLPolicyURL(client, domainID), &r.Body, nil)
	return
}

// APIACLPolicyGet retrieves details on iam identity acl for api access, by domain ID.
func APIACLPolicyGet(client *golangsdk.ServiceClient, domainID string) (r ACLResult) {
	_, r.Err = client.Get(apiACLPolicyURL(client, domainID), &r.Body, nil)
	return
}
