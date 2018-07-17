package shares

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Share contains all information associated with an OpenStack Share
type Share struct {
	// The availability zone of the share
	AvailabilityZone string `json:"availability_zone"`
	// A description of the share
	Description string `json:"description"`
	// The host name of the share
	Host string `json:"host"`
	// The UUID of the share
	ID string `json:"id"`
	// Indicates the visibility of the share
	IsPublic bool `json:"is_public"`
	// Share links for pagination
	Links []map[string]string `json:"links"`
	// Key, value -pairs of custom metadata
	Metadata map[string]string `json:"metadata"`
	// The name of the share
	Name string `json:"name"`
	// The UUID of the project to which this share belongs to
	ProjectID string `json:"project_id"`
	// The UUID of the share network
	ShareNetworkID string `json:"share_network_id"`
	// The shared file system protocol
	ShareProto string `json:"share_proto"`
	// The UUID of the share type.
	ShareType string `json:"share_type"`
	// Size of the share in GB
	Size int `json:"size"`
	// UUID of the snapshot from which to create the share
	SnapshotID string `json:"snapshot_id"`
	// The share status
	Status string `json:"status"`
	// The type of the volume
	VolumeType string `json:"volume_type"`
	// Timestamp when the share was created
	CreatedAt time.Time `json:"-"`
	//Specifies the mount location.
	ExportLocation string `json:"export_location"`
	//Lists the mount locations.
	ExportLocations []string `json:"export_locations"`
}

// AccessRight contains all information associated with an OpenStack share
// Grant Access Response
type AccessRight struct {
	// The access rule type that can be "ip", "cert" or "user".
	AccessType string `json:"access_type"`
	// The value that defines the access that can be a valid format of IP, cert or user.
	AccessTo string `json:"access_to"`
	// The access level to the share is either "rw" or "ro".
	AccessLevel string `json:"access_level"`
	// The state of the access rule
	State string `json:"state"`
	// The access rule ID.
	ID string `json:"id"`
}

// ExportLocation contains all information associated with a share export location
type ExportLocation struct {
	// The export location path that should be used for mount operation.
	Path string `json:"path"`
	// The UUID of the share instance that this export location belongs to.
	ShareInstanceID string `json:"share_instance_id"`
	// Defines purpose of an export location.
	// If set to true, then it is expected to be used for service needs
	// and by administrators only.
	// If it is set to false, then this export location can be used by end users.
	IsAdminOnly bool `json:"is_admin_only"`
	// The share export location UUID.
	ID        string `json:"id"`
	Preferred bool   `json:"preferred"`
}

// SharePage is the page returned by a pager when traversing over a
// collection of Shares.
type SharePage struct {
	pagination.LinkedPageBase
}

// Extract will get the GrantAccess object from the commonResult
func (r GrantAccessResult) ExtractAccess() (*AccessRight, error) {
	var s struct {
		AccessRight *AccessRight `json:"access"`
	}
	err := r.ExtractInto(&s)
	return s.AccessRight, err
}

// Extract will get a slice of AccessRight objects from the AccessRightsResult
func (r AccessRightsResult) ExtractAccessRights() ([]AccessRight, error) {
	var s struct {
		AccessRights []AccessRight `json:"access_list"`
	}
	err := r.ExtractInto(&s)
	return s.AccessRights, err
}

// Extract will get the Share object from the commonResult
func (r commonResult) Extract() (*Share, error) {
	var s struct {
		Share *Share `json:"share"`
	}
	err := r.ExtractInto(&s)
	return s.Share, err
}

// ExtractShares accepts a Page struct, specifically a SharePage struct,
// and extracts the elements into a slice of share structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractShares(r pagination.Page) ([]Share, error) {
	var s struct {
		ListedShares []Share `json:"shares"`
	}
	err := (r.(SharePage)).ExtractInto(&s)
	return s.ListedShares, err
}

// Extract will get the Export Locations from the commonResult
func (r GetExportLocationsResult) ExtractExportLocations() ([]ExportLocation, error) {
	var s struct {
		ExportLocations []ExportLocation `json:"export_locations"`
	}
	err := r.ExtractInto(&s)
	return s.ExportLocations, err
}

func (r *Share) UnmarshalJSON(b []byte) error {
	type tmp Share
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Share(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)

	return nil
}

// IsEmpty returns true if a ListResult contains no Shares.
func (r SharePage) IsEmpty() (bool, error) {
	shares, err := ExtractShares(r)
	return len(shares) == 0, err
}

// NextPageURL is invoked when a paginated collection of shares has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SharePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"shares_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// GrantAccessResult contains the result body and error from an GrantAccess request.
type GrantAccessResult struct {
	commonResult
}

// AccessRightsResult contains the result body and error from a AccessRight request.
type AccessRightsResult struct {
	golangsdk.Result
}

//DeleteAccessResult contains the response body from DeleteAccess rights
type DeleteAccessResult struct {
	golangsdk.Result
}

//GetExportLocationsResult contains the response body from GetExportLocations
type GetExportLocationsResult struct {
	golangsdk.Result
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}

// UpdateResult contains the response body and error from a update request.
type UpdateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// ExpandResult contains the response body and error from a Expand request.
type ExpandResult struct {
	golangsdk.ErrResult
}

// ShrinkResult contains the response body and error from a Shrink request.
type ShrinkResult struct {
	golangsdk.ErrResult
}
