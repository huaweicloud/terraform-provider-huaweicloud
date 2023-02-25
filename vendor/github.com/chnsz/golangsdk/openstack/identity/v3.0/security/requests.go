package security

import (
	"github.com/chnsz/golangsdk"
)

// PasswordPolicyOpts provides options used to update the account password policy
type PasswordPolicyOpts struct {
	// Minimum number of character types that a password must contain. Value range: 2–4.
	MinCharCombination *int `json:"password_char_combination,omitempty"`
	// Minimum number of characters that a password must contain. Value range: 6–32.
	MinPasswordLength *int `json:"minimum_password_length,omitempty"`
	// Maximum number of times that a character is allowed to consecutively present in a password.
	// Value range: 0–32.
	MaxConsecutiveIdenticalChars *int `json:"maximum_consecutive_identical_chars,omitempty"`
	// Number of previously used passwords that are not allowed. Value range: 0–10.
	RecentPasswordsDisallowedCount *int `json:"number_of_recent_passwords_disallowed,omitempty"`
	// Password validity period (days). Value range: 0–180. Value 0 indicates that this requirement does not apply.
	PasswordValidityPeriod *int `json:"password_validity_period,omitempty"`
	// Minimum period (minutes) after which users are allowed to make a password change.
	// Value range: 0–1440.
	MinPasswordAge *int `json:"minimum_password_age,omitempty"`
	// Indicates whether the password can be the username or the username spelled backwards.
	PasswordNotUsernameOrInvert *bool `json:"password_not_username_or_invert,omitempty"`
}

// UpdatePasswordPolicy can update the account password policy
func UpdatePasswordPolicy(client *golangsdk.ServiceClient, opts *PasswordPolicyOpts, domainID string) (*PasswordPolicy, error) {
	b, err := golangsdk.BuildRequestBody(opts, "password_policy")
	if err != nil {
		return nil, err
	}

	var resp PasswordPolicyResp
	_, err = client.Put(passwordPolicyURL(client, domainID), &b, &resp, nil)
	return &resp.Body, err
}

// GetPasswordPolicy retrieves details of the account password policy
func GetPasswordPolicy(client *golangsdk.ServiceClient, domainID string) (*PasswordPolicy, error) {
	var resp PasswordPolicyResp
	_, err := client.Get(passwordPolicyURL(client, domainID), &resp, nil)
	return &resp.Body, err
}
