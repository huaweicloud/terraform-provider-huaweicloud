package signs

import "github.com/chnsz/golangsdk/pagination"

// Signature is the structure that represents the signature detail.
type Signature struct {
	// The signature ID.
	ID string `json:"id"`
	// Signature key name. It can contain letters, digits, and underscores(_) and must start with a letter.
	Name string `json:"name"`
	// Signature key type.
	// + hmac
	// + basic
	// + public_key
	// + aes
	SignType string `json:"sign_type"`
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
	SignKey string `json:"sign_key"`
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
	SignSecret string `json:"sign_secret"`
	// Signature algorithm. Specify a signature algorithm only when using an AES signature key. By default, no algorithm
	// is used.
	// + aes-128-cfb
	// + aes-256-cfb
	SignAlgorithm string `json:"sign_algorithm"`
	// The creation time.
	CreatedAt string `json:"created_at"`
	// The latest update time.
	UpdatedAt string `json:"updated_at"`
	// Number of bound APIs.
	BindNum int `json:"bind_num"`
	// Number of custom backends bound.
	LdapiBindNum int `json:"ldapi_bind_num"`
}

// SignaturePage is a single page maximum result representing a query by offset page.
type SignaturePage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a SignaturePage struct is empty.
func (b SignaturePage) IsEmpty() (bool, error) {
	arr, err := ExtractSignatures(b)
	return len(arr) == 0, err
}

// ExtractSignatures is a method to extract the list of signatures.
func ExtractSignatures(r pagination.Page) ([]Signature, error) {
	var s []Signature
	err := r.(SignaturePage).Result.ExtractIntoSlicePtr(&s, "signs")
	return s, err
}

// BindResp is the structure that represents the API response of the signature binding.
type BindResp struct {
	// The published APIs of the binding relationship.
	Bindings []SignBindApiInfo `json:"bindings"`
}

// BindPage is a single page maximum result representing a query by offset page.
type BindPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a BindPage struct is empty.
func (b BindPage) IsEmpty() (bool, error) {
	arr, err := ExtractBindInfos(b)
	return len(arr) == 0, err
}

// ExtractBindInfos is a method to extract the list of binding details for signature.
func ExtractBindInfos(r pagination.Page) ([]SignBindApiInfo, error) {
	var s []SignBindApiInfo
	err := r.(BindPage).Result.ExtractIntoSlicePtr(&s, "bindings")
	return s, err
}

// SignBindApiInfo is the structure that represents the binding details.
type SignBindApiInfo struct {
	// API publish record ID.
	PublishId string `json:"publish_id"`
	// The API ID.
	ID string `json:"api_id"`
	// Group name to which the API belongs.
	GroupName string `json:"group_name"`
	// The time when the API and the signature were bound.
	BoundAt string `json:"binding_time"`
	// The environment ID where the API is published.
	EnvId string `json:"env_id"`
	// The name of the environment published by the API.
	EnvName string `json:"env_name"`
	// The API type.
	Type int `json:"api_type"`
	// The API Name.
	Name string `json:"api_name"`
	// The binding ID.
	BindId string `json:"id"`
	// The API type.
	Description string `json:"api_remark"`
	// Signature ID.
	SignId string `json:"sign_id"`
	// Signature name.
	SignName string `json:"sign_name"`
	// Signature key.
	SignKey string `json:"sign_key"`
	// Signature secret.
	SignSecret string `json:"sign_secret"`
	// Signature type.
	SignType string `json:"sign_type"`
}
