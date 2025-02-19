package users

// RequestResp is the structure that represents the API response of user methods request.
type RequestResp struct {
	// User ID.
	ID string `json:"id"`
}

// CreateResp is the structure that represents the API response of Create method request.
type CreateResp struct {
	RequestResp
}

// GetResp is the structure that represents the API response of Get method request.
type GetResp struct {
	UserDetail UserDetail `json:"user_detail"`
}

// UserDetail is an object to specified the user details.
type UserDetail struct {
	// User ID.
	ID string `json:"id"`
	// User name.
	Name string `json:"user_name"`
	// The activation mode of the user.
	// + USER_ACTIVATE: Activated by the user.
	// + ADMIN_ACTIVATE: Activated by the administator.
	ActiveType string `json:"active_type"`
	// User email.
	Email string `json:"user_email"`
	// Mobile number of the user.
	Phone string `json:"user_phone"`
	// User description.
	Description string `json:"description"`
	// User SID.
	SID string `json:"object_sid"`
	// Login account name (pre Windows).
	SamAccountName string `json:"sam_account_name"`
	// Login user name.
	UserPrincipalName string `json:"user_principal_name"`
	// Full name.
	FullName string `json:"full_name"`
	// Unique location of user on the domain tree.
	DistinguishedName string `json:"distinguished_name"`
	// Account type.
	AccountType int `json:"account_type"`
	// The character corresponding to the creation time and UTC time milliseconds.
	CreatedAt string `json:"when_created"`
	// UTC time corresponding to the last day of the account validity, in milliseconds.
	AccountExpires int `json:"account_expires"`
	// Whether the user is expired.
	UserExpired bool `json:"user_expired"`
	// Whether the account is locked.
	Locked bool `json:"locked"`
	// Whether to allow password modification.
	EnableChangePassword bool `json:"enable_change_password"`
	// Whether the password will never expires.
	PasswordNeverExpires bool `json:"password_never_expired"`
	// Whether the next login requires a password reset.
	NextLoginChangePassword bool `json:"next_login_change_password"`
	// Whether the account is disabled.
	Disabled bool `json:"disabled"`
	// Group name list.
	GroupNames []string `json:"group_names"`
	// The number of desktops the user has.
	TotalDesktops int `json:"total_desktops"`
}

// QueryResp is the structure that represents the API response of List method request.
type QueryResp struct {
	// Total number of users.
	TotalCount int `json:"total_count"`
	// User list.
	Users []User `json:"users"`
}

// User is an object to specified the some user information.
type User struct {
	// User ID.
	ID string `json:"id"`
	// User name.
	Name string `json:"user_name"`
	// User email.
	Email string `json:"user_email"`
	// The number of desktops the user has.
	TotalDesktops int `json:"total_desktops"`
	// Account expiration time, 0 means it will never expire.
	AccountExpires string `json:"account_expires"`
	// Whether the account has expired.
	AccountExpired bool `json:"account_expired"`
	// Whether the password will never expires.
	PasswordNeverExpires bool `json:"password_never_expired"`
	// Whether to allow password modification.
	EnableChangePassword bool `json:"enable_change_password"`
	// Whether the next login requires a password reset.
	NextLoginChangePassword bool `json:"next_login_change_password"`
	// User description.
	Description string `json:"description"`
	// Whether the account is locked.
	Locked bool `json:"locked"`
	// Whether the account is disabled.
	Disabled bool `json:"disabled"`
}

// UpdateResp is the structure that represents the API response of Update method request.
type UpdateResp struct {
	RequestResp
}
