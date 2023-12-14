package groups

import "github.com/chnsz/golangsdk/pagination"

// UserGroup is the structure that represents the user group detail.
type UserGroup struct {
	// The ID of user group.
	ID string `json:"id"`
	// The name of user group.
	Name string `json:"name"`
	// The type of user group.
	Type string `json:"platform_type"`
	// The description of user group.
	Description string `json:"description"`
	// The creation time of user group.
	CreatedAt string `json:"create_time"`
	// The number of users in the user group.
	UserQuantity int `json:"user_quantity"`
	// The domain ID of the user group.
	RealmID string `json:"realm_id"`
	// The specific name of user group.
	GroupDN string `json:"group_dn"`
	// The domain name of the user group.
	Domain string `json:"domain"`
	// The sID of user group.
	SID string `json:"sid"`
	// The number of desktops owned by all users in the user group.
	TotalDesktops int `json:"total_desktops"`
}

// listUserResp is the structure that represents the API response of ListUser method request.
type listUserResp struct {
	// Total number of users.
	TotalCount int `json:"total_count"`
	// User list.
	Users []User `json:"users"`
}

// User is the structure that represents the user detail under a user group.
type User struct {
	// The ID of user.
	ID string `json:"id"`
	// The name of user.
	Name string `json:"user_name"`
	// The email of user.
	Email string `json:"user_email"`
	// The phone of user.
	Phone string `json:"user_phone"`
	// The description of user .
	Description string `json:"description"`
	// The number of desktops the user has.
	TotalDesktops int `json:"total_desktops"`
}

// UserGroupPage represents the response pages of the List method.
type UserGroupPage struct {
	pagination.OffsetPageBase
}

// IsEmpty method checks whether the current user group page is empty.
func (b UserGroupPage) IsEmpty() (bool, error) {
	arr, err := ExtractUserGroupPages(b)
	return len(arr) == 0, err
}

// ExtractUserGroupPages is a method to extract the list of user group details.
func ExtractUserGroupPages(r pagination.Page) ([]UserGroup, error) {
	var s []UserGroup
	err := r.(UserGroupPage).Result.ExtractIntoSlicePtr(&s, "user_groups")
	return s, err
}
