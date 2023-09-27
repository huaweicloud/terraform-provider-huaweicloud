package custom

import "github.com/chnsz/golangsdk/pagination"

// Channel is the structure that represents the channel detail.
type Channel struct {
	// The ID of the channel.
	ID string `json:"id"`
	// The name of the channel.
	Name string `json:"name"`
	// The description of the channel.
	Description string `json:"description"`
	// The type of the channel provider.
	// + OFFICIAL: official cloud service channel.
	// + CUSTOM: the user-defined channel.
	// + PARTNER: partner channel.
	ProviderType string `json:"provider_type"`
	// The ID of the enterprise project to which the custom channel belongs.
	EnterpriseProjectId string `json:"eps_id"`
	// The creation time, in UTC format.
	CreatedTime string `json:"created_time"`
	// The update time, in UTC format.
	UpdatedTime string `json:"updated_time"`
	// The cross-account policy configuration.
	Policy CrossAccountPolicy `json:"policy"`
}

// ChannelPage is a single page maximum result representing a query by offset page.
type ChannelPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ChannelPage struct is empty.
func (b ChannelPage) IsEmpty() (bool, error) {
	arr, err := ExtractChannels(b)
	return len(arr) == 0, err
}

// ExtractChannels is a method to extract the list of custom channels.
func ExtractChannels(r pagination.Page) ([]Channel, error) {
	var s []Channel
	err := r.(ChannelPage).Result.ExtractIntoSlicePtr(&s, "items")
	return s, err
}
