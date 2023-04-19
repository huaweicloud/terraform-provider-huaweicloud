package signs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create a new signature key.
type CreateOpts struct {
	// The instnace ID to which the signature belongs.
	InstanceId string `json:"-" required:"true"`
	// Signature key name. It can contain letters, digits, and underscores(_) and must start with a letter.
	Name string `json:"name" required:"true"`
	// Signature key type.
	// + hmac
	// + basic
	// + public_key
	// + aes
	SignType string `json:"sign_type,omitempty"`
	// Signature key.
	// + For 'hmac' type: The value contains 8 to 32 characters, including letters, digits, underscores (_), and
	//   hyphens (-). It must start with a letter or digit. If not specified, a key is automatically generated.
	// + For 'basic' type: The value contains 4 to 32 characters, including letters, digits, underscores (_), and
	//   hyphens (-). It must start with a letter. If not specified, a key is automatically generated.
	// + For 'public_key' type: The value contains 8 to 512 characters, including letters, digits, and special
	//   characters (_-+/=). It must start with a letter, digit, plus sign (+), or slash (/). If not specified, a key
	//   is automatically generated.
	// + For 'aes' type: The value contains 16 characters if the aes-128-cfb algorithm is used, or 32 characters if the
	//   aes-256-cfb algorithm is used. Letters, digits, and special characters (_-!@#$%+/=) are allowed. It must start
	//   with a letter, digit, plus sign (+), or slash (/). If not specified, a key is automatically generated.
	SignKey string `json:"sign_key,omitempty"`
	// Signature secret.
	// + For 'hmac' type: The value contains 16 to 64 characters. Letters, digits, and special characters (_-!@#$%) are
	//   allowed. It must start with a letter or digit. If not specified, a value is automatically generated.
	// + For 'basic' type: The value contains 8 to 64 characters. Letters, digits, and special characters (_-!@#$%) are
	//   allowed. It must start with a letter or digit. If not specified, a value is automatically generated.
	// + For 'public_key' type: The value contains 15 to 2048 characters, including letters, digits, and special
	//   characters (_-!@#$%+/=). It must start with a letter, digit, plus sign (+), or slash (/). If not specified, a
	//   value is automatically generated.
	// + For 'aes' type: The value contains 16 characters, including letters, digits, and special
	//   characters (_-!@#$%+/=). It must start with a letter, digit, plus sign (+), or slash (/). If not specified, a
	//   value is automatically generated.
	SignSecret string `json:"sign_secret,omitempty"`
	// Signature algorithm. Specify a signature algorithm only when using an AES signature key. By default, no algorithm
	// is used.
	// + aes-128-cfb
	// + aes-256-cfb
	SignAlgorithm string `json:"sign_algorithm,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new signature key using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Signature, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Signature
	_, err = c.Post(rootURL(c, opts.InstanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts is the structure used to querying signature list.
type ListOpts struct {
	// The instnace ID to which the signature belongs.
	InstanceId string `json:"-" required:"true"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// The signature ID.
	ID string `q:"id"`
	// The signature name.
	Name string `q:"name"`
	// Parameter name (only 'name' is supported) for exact matching.
	PreciseSearch string `q:"precise_search"`
}

// List is a method to obtain the signature list under a specified instance using given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Signature, error) {
	url := rootURL(c, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := SignaturePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractSignatures(pages)
}

// UpdateOpts is the structure used to update the signature key.
type UpdateOpts struct {
	// The instnace ID to which the signature belongs.
	InstanceId string `json:"-" required:"true"`
	// The signature ID.
	SignatureId string `json:"-" required:"true"`
	// Signature key name. It can contain letters, digits, and underscores(_) and must start with a letter.
	Name string `json:"name" required:"true"`
	// Signature key type.
	// + hmac
	// + basic
	// + public_key
	// + aes
	SignType string `json:"sign_type,omitempty"`
	// Signature key.
	// + For 'hmac' type: The value contains 8 to 32 characters, including letters, digits, underscores (_), and
	//   hyphens (-). It must start with a letter or digit. If not specified, a key is automatically generated.
	// + For 'basic' type: The value contains 4 to 32 characters, including letters, digits, underscores (_), and
	//   hyphens (-). It must start with a letter. If not specified, a key is automatically generated.
	// + For 'public_key' type: The value contains 8 to 512 characters, including letters, digits, and special
	//   characters (_-+/=). It must start with a letter, digit, plus sign (+), or slash (/). If not specified, a key
	//   is automatically generated.
	// + For 'aes' type: The value contains 16 characters if the aes-128-cfb algorithm is used, or 32 characters if the
	//   aes-256-cfb algorithm is used. Letters, digits, and special characters (_-!@#$%+/=) are allowed. It must start
	//   with a letter, digit, plus sign (+), or slash (/). If not specified, a key is automatically generated.
	SignKey string `json:"sign_key,omitempty"`
	// Signature secret.
	// + For 'hmac' type: The value contains 16 to 64 characters. Letters, digits, and special characters (_-!@#$%) are
	//   allowed. It must start with a letter or digit. If not specified, a value is automatically generated.
	// + For 'basic' type: The value contains 8 to 64 characters. Letters, digits, and special characters (_-!@#$%) are
	//   allowed. It must start with a letter or digit. If not specified, a value is automatically generated.
	// + For 'public_key' type: The value contains 15 to 2048 characters, including letters, digits, and special
	//   characters (_-!@#$%+/=). It must start with a letter, digit, plus sign (+), or slash (/). If not specified, a
	//   value is automatically generated.
	// + For 'aes' type: The value contains 16 characters, including letters, digits, and special
	//   characters (_-!@#$%+/=). It must start with a letter, digit, plus sign (+), or slash (/). If not specified, a
	//   value is automatically generated.
	SignSecret string `json:"sign_secret,omitempty"`
	// Signature algorithm. Specify a signature algorithm only when using an AES signature key. By default, no algorithm
	// is used.
	// + aes-128-cfb
	// + aes-256-cfb
	SignAlgorithm string `json:"sign_algorithm,omitempty"`
}

// Update is a method used to update a signature using given parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) (*Signature, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Signature
	_, err = c.Put(resourceURL(c, opts.InstanceId, opts.SignatureId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove the specified signature using its ID and related dedicated instance ID.
func Delete(c *golangsdk.ServiceClient, instanceId, signatureId string) error {
	_, err := c.Delete(resourceURL(c, instanceId, signatureId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
