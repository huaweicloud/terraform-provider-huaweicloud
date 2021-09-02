package credentials

import "github.com/chnsz/golangsdk"

type Credential struct {
	// IAM user ID
	UserID string `json:"user_id"`

	// Description of the access key
	Description string `json:"description"`

	// AK
	AccessKey string `json:"access"`

	// SK, returned only during creation
	SecretKey string `json:"secret,omitempty"`

	// Status of the access key, active/inactive
	Status string `json:"status"`

	// Time when the access key was created
	CreateTime string `json:"create_time"`

	// Time when the access key was last used
	LastUseTime string `json:"last_use_time,omitempty"`
}

type credentialResult struct {
	golangsdk.Result
}

// CreateResult is the response of a Create operations. Call its Extract method to
// interpret it as a Credential.
type CreateResult struct {
	credentialResult
}

// Extract provides access to the Credential returned by the Get and
// Create functions.
func (r credentialResult) Extract() (*Credential, error) {
	var s struct {
		Credential *Credential `json:"credential"`
	}
	err := r.ExtractInto(&s)
	return s.Credential, err
}

// GetResult is the response of a Get operations. Call its Extract method to
// interpret it as a Credential.
type GetResult struct {
	credentialResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Credential.
type UpdateResult struct {
	credentialResult
}

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]Credential, error) {
	var a struct {
		Instances []Credential `json:"credentials"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Instances, err
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
