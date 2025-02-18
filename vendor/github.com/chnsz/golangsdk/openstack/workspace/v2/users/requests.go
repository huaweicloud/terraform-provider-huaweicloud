package users

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure required by the Create method to create a new user.
type CreateOpts struct {
	// User name.
	Name string `json:"user_name" required:"true"`
	// The activation mode of the user. Defaults to USER_ACTIVATE.
	// + USER_ACTIVATE: Activated by the user.
	// + ADMIN_ACTIVATE: Activated by the administator.
	ActiveType string `json:"active_type,omitempty"`
	// User email. The value can contain from 1 to 64 characters.
	Email string `json:"user_email,omitempty"`
	// Mobile number of the user. At least one of email and phone number must be provided.
	Phone string `json:"user_phone,omitempty"`
	// Initial passowrd of the user. The parameter is required for the administator activation mode.
	Password string `json:"password,omitempty"`
	// The expires time of Workspace user. The format is "yyyy-MM-ddTHH:mm:ss.000 Z".
	// 0 means it will never expire.
	AccountExpires string `json:"account_expires,omitempty"`
	// Whether to allow password modification.
	EnableChangePassword *bool `json:"enable_change_password,omitempty"`
	// Whether the next login requires a password reset.
	NextLoginChangePassword *bool `json:"next_login_change_password,omitempty"`
	// Group ID list.
	GroupIds []string `json:"group_ids,omitempty"`
	// User description.
	Description string `json:"description,omitempty"`
	// Alias name.
	Alias string `json:"alias_name,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a user using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain the user detail by its ID.
func Get(c *golangsdk.ServiceClient, userId string) (*UserDetail, error) {
	var r GetResp
	_, err := c.Get(resourceURL(c, userId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.UserDetail, err
}

// ListOpts is the structure required by the List method to query user list.
type ListOpts struct {
	// User name.
	Name string `q:"user_name"`
	// User description, support fuzzing match.
	Description string `q:"description"`
	// Number of records to be queried.
	// If omited, return all user's information.
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
}

// List is a method to query the job details using given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) (*QueryResp, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r QueryResp
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts is the structure required by the Update method to change user information.
type UpdateOpts struct {
	// The activation mode of the user.
	// + USER_ACTIVATE: Activated by the user.
	// + ADMIN_ACTIVATE: Activated by the administator.
	ActiveType string `json:"active_type,omitempty"`
	// User email. The value can contain from 1 to 64 characters.
	Email *string `json:"user_email,omitempty"`
	// User description.
	Description *string `json:"description,omitempty"`
	// Mobile number of the user.
	Phone *string `json:"user_phone,omitempty"`
	// The expires time of Workspace user. The format is "yyyy-MM-ddTHH:mm:ss.000Z".
	// 0 means it will never expire.
	AccountExpires string `json:"account_expires,omitempty"`
	// Whether to allow password modification.
	EnableChangePassword *bool `json:"enable_change_password,omitempty"`
	// Whether the next login requires a password reset.
	NextLoginChangePassword *bool `json:"next_login_change_password,omitempty"`
	// Whether the password will never expires.
	PasswordNeverExpires *bool `json:"password_never_expired,omitempty"`
	// Whether the account is disabled.
	Disabled *bool `json:"disabled,omitempty"`
}

// Update is a method to change user informaion using givin parameters.
func Update(c *golangsdk.ServiceClient, userId string, opts UpdateOpts) (*UpdateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r UpdateResp
	_, err = c.Put(resourceURL(c, userId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove an existing user using given parameters.
func Delete(c *golangsdk.ServiceClient, userId string) error {
	_, err := c.Delete(resourceURL(c, userId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
