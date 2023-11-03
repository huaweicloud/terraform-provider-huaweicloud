package endpoints

import "github.com/chnsz/golangsdk/pagination"

// CreateResp is the structure that represents the response to adding permission records.
type createResp struct {
	// The permission list of endpoint service.
	Permissions []string `json:"permissions"`
}

// EndpointPermission is the structure that represents the permission detail.
type EndpointPermission struct {
	// The permission ID.
	ID string `json:"id"`
	// The permission of endpoint service.
	Permission string `json:"permission"`
	// The creation time, in UTC format.
	CreatedAt string `json:"created_at"`
}

// PermissionPage is a single page maximum result representing a query by offset page.
type PermissionPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a PermissionPage struct is empty.
func (b PermissionPage) IsEmpty() (bool, error) {
	arr, err := ExtractPermissions(b)
	return len(arr) == 0, err
}

// ExtractPermissions is a method to extract the list of permissions for specified endpoint service.
func ExtractPermissions(r pagination.Page) ([]EndpointPermission, error) {
	var s []EndpointPermission
	err := r.(PermissionPage).Result.ExtractIntoSlicePtr(&s, "permissions")
	return s, err
}
