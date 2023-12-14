package groups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to create a new user group.
type CreateOpts struct {
	// The name of user group.
	Name string `json:"group_name" required:"true"`
	// The type of user group. The valid types are as following:
	//   AD: AD domain user group
	//   LOCAL: Local liteAs user group
	Type string `json:"platform_type" required:"true"`
	// The description of user group.
	Description string `json:"description,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a user group using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(rootURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// ListOpts is the structure required by the List method to query user group list.
type ListOpts struct {
	// Search keywords used to match user groups. For example, fuzzy queries based on user group names.
	Keyword string `q:"keyword"`
	// Number of records to be queried.
	// Valid range is 0-100.
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
}

// List is a method to query the user group details using given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]UserGroup, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := UserGroupPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})

	pager.Headers = requestOpts.MoreHeaders
	pages, err := pager.AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractUserGroupPages(pages)
}

// UpdateOpts is the structure required by the Update method to change user group information.
type UpdateOpts struct {
	// The ID of user group.
	GroupID string `json:"-" required:"true"`
	// The name of user group.
	Name string `json:"group_name,omitempty"`
	// The description of user group.
	Description *string `json:"description,omitempty"`
}

// Update is a method to change user group information using given parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(resourceURL(c, opts.GroupID), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// Delete is a method to remove an existing user group using given parameters.
func Delete(c *golangsdk.ServiceClient, groupId string) error {
	_, err := c.Delete(resourceURL(c, groupId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// ListUserOpts is the structure required by the ListUser method to query all user information under a user group.
type ListUserOpts struct {
	// The ID of user.
	Name string `q:"user_name"`
	// The description of user.
	Description string `q:"description"`
	// The activation type of user. The valid types are as following:
	//   USER_ACTIVATE: User activation
	//   ADMIN_ACTIVATE: Administrator activation
	Type string `q:"active_type"`
	// Number of records to be queried.
	// Valid range is 0-2000.
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
}

// ListUser is a method to query all user information under a user group using given parameters.
func ListUser(c *golangsdk.ServiceClient, groupId string, opts ListUserOpts) ([]User, error) {
	url := userURL(c, groupId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r listUserResp
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Users, err
}

// ActionOpts is the structure required by the DoAction method of add users to or remove users from user group.
type ActionOpts struct {
	// The user ID to be added or removed from the user group.
	UserIDs []string `json:"user_ids" required:"true"`
	// The operation type. The valid types are as following:
	//   ADD: Add user
	//   DELETE: Delete user
	Type string `json:"op_type" required:"true"`
}

// DoAction is a method of add users to or remove users from user group using given parameters.
func DoAction(client *golangsdk.ServiceClient, groupId string, opts ActionOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}
	_, err = client.Post(actionURL(client, groupId), &b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return err
}
