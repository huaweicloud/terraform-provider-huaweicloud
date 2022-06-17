package rotation

import (
	"github.com/chnsz/golangsdk"
)

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	golangsdk.Result
}

// Rotation is the response body of Get method
type Rotation struct {
	// Key rotation status. The default value is false, indicating that key rotation is disabled.
	Enabled bool `json:"key_rotation_enabled"`
	// Rotation interval. The value is an integer in the range 30 to 365.
	Interval int `json:"rotation_interval"`
	// Last key rotation time. The timestamp indicates the total microseconds past the start of the epoch date (January 1, 1970).
	LastRotationTime string `json:"last_rotation_time"`
	// Number of key rotations.
	NumberOfRotations int `json:"number_of_rotations"`
}

func (r GetResult) Extract() (Rotation, error) {
	var s Rotation
	err := r.ExtractInto(&s)
	return s, err
}
