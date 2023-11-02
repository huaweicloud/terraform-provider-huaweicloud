package custom

import "github.com/chnsz/golangsdk/pagination"

// Source is the structure that represents the event source detail.
type Source struct {
	// The ID of the event source.
	ID string `json:"id"`
	// The name of the event source.
	Name string `json:"name"`
	// The description of the event source.
	Description string `json:"description"`
	// The type of the event source provider.
	// + OFFICIAL: official cloud service event source.
	// + CUSTOM: the user-defined event source.
	// + PARTNER: partner event source.
	ProviderType string `json:"provider_type"`
	// The list of event types provided by the event source.
	// Only the official cloud service event source provides event types.
	EventTypes []EventType `json:"event_types"`
	// The creation time, in UTC format.
	CreatedTime string `json:"created_time"`
	// The update time, in UTC format.
	UpdatedTime string `json:"updated_time"`
	// The ID of the event channel to which the event source belongs.
	ChannelId string `json:"channel_id"`
	// The name of the event channel to which the event source belongs.
	ChannelName string `json:"channel_name"`
	// The type of the event source.
	// + APPLICATION (default)
	// + RABBITMQ
	// + ROCKETMQ
	Type string `json:"type"`
	// The configuration detail of the event source, in JSON format.
	Detail interface{} `json:"detail"`
	// The status of the event source
	Status string `json:"status"`
}

// EventType is the structure that represents the event type that provided by the event source.
type EventType struct {
	// The name of the event type.
	Name string `json:"name"`
	// The description of the event type.
	Description string `json:"description"`
}

// SourcePage is a single page maximum result representing a query by offset page.
type SourcePage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a SourcePage struct is empty.
func (b SourcePage) IsEmpty() (bool, error) {
	arr, err := ExtractSources(b)
	return len(arr) == 0, err
}

// ExtractSources is a method to extract the list of custom event sources.
func ExtractSources(r pagination.Page) ([]Source, error) {
	var s []Source
	err := r.(SourcePage).Result.ExtractIntoSlicePtr(&s, "items")
	return s, err
}
