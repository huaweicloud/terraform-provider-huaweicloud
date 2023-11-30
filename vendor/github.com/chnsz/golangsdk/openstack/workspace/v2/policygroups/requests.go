package policygroups

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure that represents the policy group configuration.
type CreateOpts struct {
	// Policy group name.
	PolicyGroupName string `json:"policy_group_name" required:"true"`
	// Policy group ID.
	PolicyGroupId string `json:"policy_group_id,omitempty"`
	// Priority.
	Priority int `json:"priority,omitempty"`
	// Policy group description.
	Description string `json:"description,omitempty"`
	// The update time of the policy group.
	UpdateTime string `json:"update_time,omitempty"`
	// List of target objects.
	Targets []Target `json:"targets,omitempty"`
	// Access policy configuration.
	Policies Policy `json:"policies,omitempty"`
}

// Target is the structure that represents the access target configuration.
type Target struct {
	// Target ID.
	// If the target type is 'INSTANCE', the ID means the SID of the desktop.
	// If the target type is 'USER', the ID means the user ID.
	// If the target type is 'USERGROUP', the ID means the user group ID.
	// If the target type is 'CLIENTIP', the ID means the terminal IP address.
	// If the target type is 'OU', the ID means the OUID.
	// If the target type is 'ALL', the ID fixed with string 'default-apply-all-targets'.
	TargetId string `json:"target_id,omitempty"`
	// Target name.
	// If the target type is 'INSTANCE', the ID means the desktop name.
	// If the target type is 'USER', the ID means the user name.
	// If the target type is 'USERGROUP', the ID means the user group name.
	// If the target type is 'CLIENTIP', the ID means the terminal IP address.
	// If the target type is 'OU', the ID means the OU name.
	// If the target type is 'ALL', the ID fixed with string 'All-Targets'.
	TargetName string `json:"target_name,omitempty"`
	// Target type.
	// + INSTANCE: Desktop.
	// + USER: User.
	// + USERGROUP: User group.
	// + CLIENTIP: Terminal IP address.
	// + OU: Organization unit.
	// + ALL: All desktops.
	TargetType string `json:"target_type,omitempty"`
}

// Policy is the structure that represents the access policy detail.
type Policy struct {
	// Access control.
	AccessControl AccessControl `json:"access_control,omitempty"`
}

// AccessControl is the structure that represents the access policy detail.
type AccessControl struct {
	// IP access control.
	// It consists of multiple groups of IP addresses and network masks, separated by ';',
	// and spliced together by '|' between IP addresses and network masks.
	// Such as: "IP|mask;IP|mask;IP|mask"
	IpAccessControl string `json:"ip_access_control,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a policy group using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "policy_group")
	if err != nil {
		return "", err
	}

	var r struct {
		ID string `json:"id"`
	}
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.ID, err
}

// Get is a method to obtain the policy group detail by its ID.
func Get(c *golangsdk.ServiceClient, groupId string) (*PolicyGroup, error) {
	var r struct {
		PolicyGroup PolicyGroup `json:"policy_group"`
	}
	_, err := c.Get(resourceURL(c, groupId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.PolicyGroup, err
}

// UpdateOpts is the structure that used to modify the policy group configuration.
type UpdateOpts struct {
	// Policy group ID.
	PolicyGroupId string `json:"-" required:"true"`
	// Policy group name.
	PolicyGroupName string `json:"policy_group_name" required:"true"`
	// Priority.
	Priority int `json:"priority,omitempty"`
	// Policy group description.
	Description *string `json:"description,omitempty"`
	// List of target objects.
	Targets []Target `json:"targets,omitempty"`
	// List of policy.
	Policies *Policy `json:"policies,omitempty"`
}

// Update is a method to modify the existing policy group configuration using given parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "policy_group")
	if err != nil {
		return "", err
	}

	var r struct {
		ID string `json:"id"`
	}
	_, err = c.Put(resourceURL(c, opts.PolicyGroupId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.ID, err
}

// Delete is a method to remove an existing policy group using given parameters.
func Delete(c *golangsdk.ServiceClient, groupId string) error {
	_, err := c.Delete(resourceURL(c, groupId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
