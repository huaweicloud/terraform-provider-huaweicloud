package whitelists

import (
	"github.com/huaweicloud/golangsdk"
)

// WhitelistOptsBuilder is used for creating, updating, deleting instance whitelists parameters.
// any struct providing the parameters should implement this interface
type WhitelistOptsBuilder interface {
	ToInstanceWhitelistMap() (map[string]interface{}, error)
}

// WhitelistGroupOpts is a struct that contains all the whitelist parameters.
type WhitelistGroupOpts struct {
	// the group name
	GroupName string `json:"group_name" required:"true"`
	// list of IP address or range
	IPList []string `json:"ip_list" required:"true"`
}

// WhitelistOpts is a struct that contains all the parameters.
type WhitelistOpts struct {
	// enable or disable the whitelists
	Enable *bool `json:"enable_whitelist" required:"true"`
	// list of whitelist groups
	Groups []WhitelistGroupOpts `json:"whitelist" required:"true"`
}

// ToInstanceWhitelistMap is used for type convert
func (ops WhitelistOpts) ToInstanceWhitelistMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Put an instance whitelist with given parameters.
func Put(client *golangsdk.ServiceClient, id string, ops WhitelistOptsBuilder) (r PutResult) {
	b, err := ops.ToInstanceWhitelistMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(resourceURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Get the instance whitelist groups by instance id
func Get(client *golangsdk.ServiceClient, id string) (r WhitelistResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}
