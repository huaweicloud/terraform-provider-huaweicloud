package keys

import (
	"github.com/gophercloud/gophercloud"
)

type CreateOpts struct {
	// Alias of a CMK
	KeyAlias string `json:"key_alias" required:"true"`
	// CMK description
	KeyDescription string `json:"key_description,omitempty"`
	// Region where a CMK resides
	Realm string `json:"realm,omitempty"`
	// Key policy (This parameter is an extension field.)
	KeyPolicy string `json:"key_policy,omitempty"`
	// Purpose of a CMK (The default value is Encrypt_Decrypt)
	KeyUsage string `json:"key_usage,omitempty"`
	// Type of a CMK
	KeyType string `json:"key_type,omitempty"`
	// 36-byte serial number of a request message
	Sequence string `json:"sequence,omitempty"`
}

// ListOpts holds options for getting keys. It is passed to the keys.Get
// function.
type ListOpts struct {
	// ID of a CMK
	KeyID string `json:"key_id" required:"true"`
}

type DeleteOpts struct {
	// ID of a CMK
	KeyID string `json:"key_id" required:"true"`
	// Number of days after which a CMK is scheduled to be deleted
	// (The value ranges from 7 to 1096.)
	PendingDays string `json:"pending_days" required:"true"`
	// 36-byte serial number of a request message
	Sequence string `json:"sequence,omitempty"`
}

type UpdateAliasOpts struct {
	// ID of a CMK
	KeyID string `json:"key_id" required:"true"`
	// CMK description
	KeyAlias string `json:"key_alias" required:"true"`
	// 36-byte serial number of a request message
	Sequence string `json:"sequence,omitempty"`
}

type UpdateDesOpts struct {
	// ID of a CMK
	KeyID string `json:"key_id" required:"true"`
	// CMK description
	KeyDescription string `json:"key_description" required:"true"`
	// 36-byte serial number of a request message
	Sequence string `json:"sequence,omitempty"`
}


// ToKeyCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToKeyCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToGetKeyMap formats an ListOpts structure into a request body.
func (opts ListOpts) ToGetKeyMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToKeyDeleteMap assembles a request body based on the contents of a
// DeleteOpts.
func (opts DeleteOpts) ToKeyDeleteMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToKeyUpdateAliasMap assembles a request body based on the contents of a
// UpdateAliasOpts.
func (opts UpdateAliasOpts) ToKeyUpdateAliasMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ToKeyUpdateDesMap assembles a request body based on the contents of a
// UpdateDesOpts.
func (opts UpdateDesOpts) ToKeyUpdateDesMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type CreateOptsBuilder interface {
	ToKeyCreateMap() (map[string]interface{}, error)
}

type GetOptsBuilder interface {
	ToGetKeyMap() (map[string]interface{}, error)
}

type DeleteOptsBuilder interface {
	ToKeyDeleteMap() (map[string]interface{}, error)
}

type UpdateAliasOptsBuilder interface {
	ToKeyUpdateAliasMap() (map[string]interface{}, error)
}

type UpdateDesOptsBuilder interface {
	ToKeyUpdateDesMap() (map[string]interface{}, error)
}

// Create will create a new key based on the values in CreateOpts. To ExtractKeyInfo
// the key object from the response, call the ExtractKeyInfo method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToKeyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves the key with the provided ID. To extract the key object
// from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, opts GetOptsBuilder) (r GetResult) {
	b, err := opts.ToGetKeyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(getURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will delete the existing key with the provided ID.
func Delete(client *gophercloud.ServiceClient, opts DeleteOptsBuilder) (r DeleteResult) {
	b, err := opts.ToKeyDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
		JSONResponse: &r.Body,
	})
	return
}

func UpdateAlias(client *gophercloud.ServiceClient, opts UpdateAliasOptsBuilder) (r UpdateAliasResult) {
	b, err := opts.ToKeyUpdateAliasMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateAliasURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func UpdateDes(client *gophercloud.ServiceClient, opts UpdateDesOptsBuilder) (r UpdateDesResult) {
	b, err := opts.ToKeyUpdateDesMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateDesURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
