package aliases

// Alias is the structure that represent the details of the version alias.
type Alias struct {
	// Function alias.
	Name string `json:"name"`
	// Version corresponding to the alias.
	Version string `json:"version"`
	// Description of the alias.
	Description string `json:"description"`
	// Time when the alias was last modified.
	LastModified string `json:"last_modified"`
	// URN of the alias.
	AliasUrn string `json:"alias_urn"`
	// The weights configuration of the additional version.
	AdditionalVersionWeights map[string]interface{} `json:"additional_version_weights"`
}
