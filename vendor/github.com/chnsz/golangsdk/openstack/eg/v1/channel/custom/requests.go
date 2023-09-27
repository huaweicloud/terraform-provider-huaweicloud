package custom

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create a new custom channel.
type CreateOpts struct {
	// The name of the custom channel.
	Name string `json:"name" required:"true"`
	// The description of the custom channel.
	Description string `json:"description,omitempty"`
	// The ID of the enterprise project to which the custom channel belongs.
	EnterpriseProjectId string `json:"eps_id,omitempty" q:"enterprise_project_id"`
	// Whether enable cross-account configuration.
	CrossAccount *bool `json:"cross_account,omitempty"`
	// The event policy of the cross-account.
	Policy *CrossAccountPolicy `json:"policy,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new custom channel using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*Channel, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r Channel
	_, err = client.Post(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to query an existing channel by its ID.
func Get(client *golangsdk.ServiceClient, channelId, epsId string) (*Channel, error) {
	var r Channel
	url := resourceURL(client, channelId)
	if epsId != "" {
		url = fmt.Sprintf("%s?enterprise_project_id=%s", url, epsId)
	}
	_, err := client.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts is the structure used to query channel list.
type ListOpts struct {
	// The ID of the event channel to which the channel belongs.
	ChannelId string `q:"channel_id"`
	// The ID of the enterprise project to which the custom channel belongs.
	EnterpriseProjectId string `q:"enterprise_project_id"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0.
	// The valid value ranges from 0 to 100, defaults to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	// The valid value ranges from 1 to 1000, defaults to 15.
	Limit int `q:"limit"`
	// The query sorting.
	// The default value is 'created_time:DESC'.
	Sort string `q:"sort"`
	// The type of the channel provider.
	// + OFFICIAL: official cloud service channel.
	// + CUSTOM: the user-defined channel.
	// + PARTNER: partner channel.
	ProviderType string `q:"provider_type"`
	// The name of the channel.
	Name string `q:"name"`
	// The fuzzy name of the channel.
	FuzzyName string `q:"fuzzy_name"`
}

// List is a method to query the channel list using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Channel, error) {
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ChannelPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractChannels(pages)
}

// UpdateOpts is the structure used to update an existing custom channel.
type UpdateOpts struct {
	// The ID of the channel.
	ChannelId string `json:"-" required:"true"`
	// The description of the channel.
	Description *string `json:"description,omitempty"`
	// Whether enable cross-account configuration.
	CrossAccount *bool `json:"cross_account,omitempty"`
	// The event policy of the cross-account.
	Policy *CrossAccountPolicy `json:"policy,omitempty"`
	// The ID of the enterprise project to which the custom channel belongs.
	// Notes: this parameter does not support update, but it is required for request body and query parameter.
	EnterpriseProjectId string `json:"eps_id,omitempty" q:"enterprise_project_id"`
}

// CrossAccountPolicy is the structure that represents the cross-account policy.
type CrossAccountPolicy struct {
	// The SID of the cross-account policy.
	Sid string `json:"Sid"`
	// The effect of the cross-account policy.
	// + Allow
	// + Deny
	Effect string `json:"Effect"`
	// The configuration of the IAM account.
	Principal PrincipalInfo `json:"Principal"`
	// The action of the cross-account policy, such as 'eg:channels:putEvents'.
	Action string `json:"Action"`
	// The URN of the custom channel.
	// The format is 'urn:eg:{region}:{domain_id}:channel:{channel_name}'
	// Before channel created, the channel name is empty.
	Resource string `json:"Resource"`
}

// PrincipalInfo is the structure that represents the domain ID list of the cross-account policy.
type PrincipalInfo struct {
	// The account IDs.
	IAM []string `json:"IAM"`
}

// Update is a method used to modify an existing custom channel using given parameters.
func Update(client *golangsdk.ServiceClient, opts UpdateOpts) (*Channel, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	url := resourceURL(client, opts.ChannelId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r Channel
	_, err = client.Put(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing custom channel using its ID.
func Delete(client *golangsdk.ServiceClient, channelId, epsId string) error {
	url := resourceURL(client, channelId)
	if epsId != "" {
		url = fmt.Sprintf("%s?enterprise_project_id=%s", url, epsId)
	}
	_, err := client.Delete(url, nil)
	return err
}
