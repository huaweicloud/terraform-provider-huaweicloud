package rotation

import (
	"github.com/chnsz/golangsdk"
)

type RotationOptsBuilder interface {
	ToKeyRotationMap() (map[string]interface{}, error)
}

type UpdateOptsBuilder interface {
	ToKeyRotationIntervalMap() (map[string]interface{}, error)
}

type RotationOpts struct {
	// ID of a CMK
	KeyID string `json:"key_id" required:"true"`
	// 36-byte sequence number of a request message
	Sequence string `json:"sequence,omitempty"`
}

type IntervalOpts struct {
	// ID of a CMK
	KeyID string `json:"key_id" required:"true"`
	// Rotation interval. The value is an integer in the range 30 to 365.
	// Set the interval based on how often a CMK is used.
	// If it is frequently used, set a short interval; otherwise, set a long one.
	Interval int `json:"rotation_interval" required:"true"`
	// 36-byte sequence number of a request message
	Sequence string `json:"sequence,omitempty"`
}

// ToKeyRotationMap assembles a request body based on the contents of a
// RotationOpts.
func (opts RotationOpts) ToKeyRotationMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ToKeyRotationIntervalMap assembles a request body based on the contents of a
// IntervalOpts.
func (opts IntervalOpts) ToKeyRotationIntervalMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Enable will enable rotation for a CMK based on the values in RotationOpts.
// The default rotation interval is 365 days.
// CMKs created using imported key materials and Default Master Keys do not support rotation.
func Enable(client *golangsdk.ServiceClient, opts RotationOptsBuilder) (r golangsdk.ErrResult) {
	b, err := opts.ToKeyRotationMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(enableURL(client), b, &r.Body, nil)
	return
}

// Disable will disable rotation for a CMK based on the values in RotationOpts.
func Disable(client *golangsdk.ServiceClient, opts RotationOptsBuilder) (r golangsdk.ErrResult) {
	b, err := opts.ToKeyRotationMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(disableURL(client), b, &r.Body, nil)
	return
}

// Get retrieves the key with the provided ID. To extract the key object
// from the response, call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, opts RotationOptsBuilder) (r GetResult) {
	b, err := opts.ToKeyRotationMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(getURL(client), &b, &r.Body, nil)
	return
}

// Update will change the rotation interval for a CMK
func Update(client *golangsdk.ServiceClient, opts UpdateOptsBuilder) (r golangsdk.ErrResult) {
	b, err := opts.ToKeyRotationIntervalMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(intervalURL(client), b, &r.Body, nil)
	return
}
