package subscriptions

import "github.com/chnsz/golangsdk/pagination"

// Subscription is the structure that represents the event subscription detail.
type Subscription struct {
	// The ID of the event subscription.
	ID string `json:"id"`
	// The name of the event subscription.
	Name string `json:"name"`
	// The description of the event subscription.
	Description string `json:"description"`
	// The type of the event subscription.
	// + EVENT
	// + SCHEDULED
	Type string `json:"type"`
	// The status of the event subscription.
	Status string `json:"status"`
	// The ID of the event channel to which the event subscription belongs.
	ChannelId string `json:"channel_id"`
	// The name of the event channel to which the event subscription belongs.
	ChannelName string `json:"channel_name"`
	// The list of used resources.
	Used []UsedInfo `json:"used"`
	// The list of the event sources.
	Sources []map[string]interface{} `json:"sources"`
	// The list of the event targets.
	Targets []map[string]interface{} `json:"targets"`
	// The creation time, in UTC format.
	CreatedTime string `json:"created_time"`
	// The update time, in UTC format.
	UpdatedTime string `json:"updated_time"`
}

// UsedInfo is the structure that represents the detail of the used resources.
type UsedInfo struct {
	// Associated resource ID.
	ResourceId string `json:"resource_id"`
	// The name of the tenant that owns the resource.
	Owner string `json:"owner"`
	// The description.
	Description string `json:"description"`
}

// SubscriptionPage is a single page maximum result representing a query by offset page.
type SubscriptionPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a SubscriptionPage struct is empty.
func (b SubscriptionPage) IsEmpty() (bool, error) {
	arr, err := ExtractSubscriptions(b)
	return len(arr) == 0, err
}

// ExtractSubscriptions is a method to extract the list of subscriptions.
func ExtractSubscriptions(r pagination.Page) ([]Subscription, error) {
	var s []Subscription
	err := r.(SubscriptionPage).Result.ExtractIntoSlicePtr(&s, "items")
	return s, err
}
