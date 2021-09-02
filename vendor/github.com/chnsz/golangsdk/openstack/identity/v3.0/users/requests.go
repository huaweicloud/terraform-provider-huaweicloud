package users

import (
	"github.com/chnsz/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a user.
type CreateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name" required:"true"`

	// DomainID is the ID of the domain the user belongs to.
	DomainID string `json:"domain_id" required:"true"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`

	// Email address with a maximum of 255 characters
	Email string `json:"email,omitempty"`

	// AreaCode is a country code, must be used together with Phone.
	AreaCode string `json:"areacode,omitempty"`

	// Phone is a mobile number with a maximum of 32 digits, must be used together with AreaCode.
	Phone string `json:"phone,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`

	// AccessMode is the access type for IAM user
	AccessMode string `json:"access_mode,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// PasswordReset Indicates whether password reset is required at the first login.
	// By default, password reset is true.
	PasswordReset *bool `json:"pwd_status,omitempty"`
}

// ToUserCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToUserCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Create creates a new User.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToUserUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts provides options for updating a user account.
type UpdateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name,omitempty"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`

	// Email address with a maximum of 255 characters
	Email string `json:"email,omitempty"`

	// AreaCode is a country code, must be used together with Phone.
	AreaCode string `json:"areacode,omitempty"`

	// Phone is a mobile number with a maximum of 32 digits. must be used together with AreaCode.
	Phone string `json:"phone,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`

	// AccessMode is the access type for IAM user
	AccessMode string `json:"access_mode,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// PasswordReset Indicates whether password reset is required
	PasswordReset *bool `json:"pwd_status,omitempty"`
}

// ToUserUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToUserUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update updates an existing User.
func Update(client *golangsdk.ServiceClient, userID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUserUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, userID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves details on a single user, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
