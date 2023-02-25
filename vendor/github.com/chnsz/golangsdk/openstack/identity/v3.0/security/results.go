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
