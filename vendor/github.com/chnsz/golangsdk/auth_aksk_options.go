package golangsdk

// AKSKAuthOptions presents the required information for AK/SK auth
type AKSKAuthOptions struct {
	// IdentityEndpoint specifies the HTTP endpoint that is required to work with
	// the Identity API of the appropriate version. While it's ultimately needed by
	// all of the identity services, it will often be populated by a provider-level
	// function.
	//
	// The IdentityEndpoint is typically referred to as the "auth_url" or
	// "OS_AUTH_URL" in the information provided by the cloud operator.
	IdentityEndpoint string `json:"-"`

	// region
	Region string

	// IAM user project id and name
	ProjectId   string
	ProjectName string

	// IAM account name and id
	Domain   string
	DomainID string

	// cloud service domain for BSS
	BssDomain   string
	BssDomainID string

	AccessKey     string //Access Key
	SecretKey     string //Secret key
	SecurityToken string //Security Token for temporary Access Key

	IsDerived bool // Whether to enable the derivative algorithm

	DerivedAuthServiceName string // Derivative algorithm service name. Required for derivative algorithm.

	// AgencyNmae is the name of agnecy
	AgencyName string

	// AgencyDomainName is the name of domain who created the agency
	AgencyDomainName string

	// DelegatedProject is the name of delegated project
	DelegatedProject string

	// whether using the customer catalog, defaults to false
	WithUserCatalog bool

	// SigningAlgorithm is used to select encryption algorithm
	SigningAlgorithm string
}

// GetIdentityEndpoint implements the method of AuthOptionsProvider
func (opts AKSKAuthOptions) GetIdentityEndpoint() string {
	return opts.IdentityEndpoint
}
