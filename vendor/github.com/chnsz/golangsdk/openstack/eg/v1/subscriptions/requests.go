package subscriptions

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create a new subscription.
type CreateOpts struct {
	// The ID of the event channel associated to the subscription.
	ChannelId string `json:"channel_id" required:"true"`
	// The name of the subscription.
	// The valid length is limited from `1` to `128`, only letters, digits, hyphens (-), underscores (_) and dots (.)
	// are allowed. The name must start with a letter or digit.
	Name string `json:"name" required:"true"`
	// The list of the event sources.
	Sources []interface{} `json:"sources" required:"true"`
	// The list of the event targets.
	Targets []interface{} `json:"targets" required:"true"`
	// The description of the subscription.
	// The valid length is limited from `0` to `128`.
	Description string `json:"description,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new subscription using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*Subscription, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Subscription
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to query an existing subscription by its ID.
func Get(client *golangsdk.ServiceClient, subscriptionId string) (*Subscription, error) {
	var r Subscription
	_, err := client.Get(resourceURL(client, subscriptionId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts is the structure used to query event subscription list.
type ListOpts struct {
	// The ID of the event channel associated to subscription.
	ChannelId string `q:"channel_id"`
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
	// The name of the subscription.
	Name string `q:"name"`
	// The fuzzy name of the subscription.
	FuzzyName string `q:"fuzzy_name"`
	// The fuzzy label of the subscription.
	ConnectionId string `q:"connection_id"`
}

// List is a method to query the subscription list using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Subscription, error) {
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := SubscriptionPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractSubscriptions(pages)
}

// UpdateOpts is the structure used to update an existing subscription.
type UpdateOpts struct {
	SubscriptionId string `json:"-" required:"true"`
	// The description of the subscription.
	Description *string `json:"description,omitempty"`
	// The list of the event sources.
	Sources []interface{} `json:"sources,omitempty"`
	// The list of the event targets.
	Targets []interface{} `json:"targets,omitempty"`
}

// Update is a method used to modify an existing event subscription using given parameters.
func Update(client *golangsdk.ServiceClient, opts UpdateOpts) (*Subscription, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Subscription
	_, err = client.Put(resourceURL(client, opts.SubscriptionId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing event subscription using its ID.
func Delete(client *golangsdk.ServiceClient, subscriptionId string) error {
	_, err := client.Delete(resourceURL(client, subscriptionId), nil)
	return err
}
