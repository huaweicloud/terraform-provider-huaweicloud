package providers

import "github.com/chnsz/golangsdk/pagination"

// Provider is the structure that represents the resources detail of the tags supported providers.
type Provider struct {
	// Specifies the cloud service name.
	Provider string `json:"provider"`
	// Specifies the display name of the resource. You can configure the language by setting the locale parameter.
	ProviderI18nDisplayName string `json:"provider_i18n_display_name"`
	// Specifies the resource type list.
	Resources []ResourceDetail `json:"resource_types"`
}

// ResourceDetail is the strucutre that represents the resource detail.
type ResourceDetail struct {
	// Specifies the resource type.
	ResourceType string `json:"resource_type"`
	// Specifies the display name of the resource type. You can configure the language by setting the locale parameter.
	ResourceTypeI18nDisplayName string `json:"resource_type_i18n_display_name"`
	// Specifies supported regions.
	Regions []string `json:"regions"`
	// Specifies whether the resources is a global resource.
	Global bool `json:"global"`
}

// ProviderPage is a single page maximum result representing a query by offset page.
type ProviderPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ProviderPage struct is empty.
func (b ProviderPage) IsEmpty() (bool, error) {
	arr, err := extractProviders(b)
	return len(arr) == 0, err
}

// extractProviders is a method to extract the list of tags supported providers.
func extractProviders(r pagination.Page) ([]Provider, error) {
	var s []Provider
	err := r.(ProviderPage).Result.ExtractIntoSlicePtr(&s, "providers")
	return s, err
}
