package security

type PasswordPolicy struct {
	// Minimum number of character types that a password must contain. Value range: 2–4.
	MinCharCombination int `json:"password_char_combination"`
	// Minimum number of characters that a password must contain. Value range: 6–32.
	MinPasswordLength int `json:"minimum_password_length"`
	// Maximum number of characters that a password can contain.
	MaxPasswordLength int `json:"maximum_password_length"`
	// Maximum number of times that a character is allowed to consecutively present in a password.
	MaxConsecutiveIdenticalChars int `json:"maximum_consecutive_identical_chars"`
	// Number of previously used passwords that are not allowed. Value range: 0–10.
	RecentPasswordsDisallowedCount int `json:"number_of_recent_passwords_disallowed"`
	// Password validity period (days). Value range: 0–180. Value 0 indicates that this requirement does not apply.
	PasswordValidityPeriod int `json:"password_validity_period"`
	// Minimum period (minutes) after which users are allowed to make a password change.
	MinPasswordAge int `json:"minimum_password_age"`
	// Indicates whether the password can be the username or the username spelled backwards.
	PasswordNotUsernameOrInvert bool `json:"password_not_username_or_invert"`
	// Characters that a password must contain.
	PasswordRequirements string `json:"password_requirements"`
}

type PasswordPolicyResp struct {
	Body PasswordPolicy `json:"password_policy"`
}

type ProtectPolicy struct {
	// Indicates whether to enable operation protection.
	Protection bool `json:"operation_protection"`
	// the attributes IAM users can modify.
	AllowUser AllowUserBody `json:"allow_user"`
	// Indicates whether to designate a person for verification.
	AdminCheck string `json:"admin_check"`
	// the verification method. The optional values are mobile and email.
	Scene string `json:"scene"`
	// the mobile number used for verification. Example: 0086-123456789
	Mobile string `json:"mobile"`
	// the email address used for verification. An example value is example@email.com.
	Email string `json:"email"`
}

type AllowUserBody struct {
	// Indicates whether to allow IAM users to manage access keys by themselves.
	ManageAccesskey bool `json:"manage_accesskey"`
	// Indicates whether to allow IAM users to change their email addresses.
	ManageEmail bool `json:"manage_email"`
	// Indicates whether to allow IAM users to change their mobile numbers.
	ManageMobile bool `json:"manage_mobile"`
	// Indicates whether to allow IAM users to change their passwords.
	ManagePassword bool `json:"manage_password"`
}

type ProtectPolicyResp struct {
	Body ProtectPolicy `json:"protect_policy"`
}
