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

// ProtectPolicyOpts provides options used to modify the operation protection policy.
type ProtectPolicyOpts struct {
	// Specifies whether to enable operation protection. The value can be true or false.
	Protection *bool `json:"operation_protection" required:"true"`
	// Specifies the attributes IAM users can modify.
	AllowUser *AllowUserOpts `json:"allow_user,omitempty"`
	// Specifies whether to designate a person for verification.
	// If this parameter is set to on, you need to specify the scene parameter to designate a person for verification.
	// If this parameter is set to off, the designated operator is responsible for the verification.
	AdminCheck *string `json:"admin_check,omitempty"`
	// Specifies the verification method. This parameter is mandatory when admin_check is set to on.
	// The optional values are mobile and email.
	Scene *string `json:"scene,omitempty"`
	// Specifies the mobile number used for verification. Example: 0086-123456789
	Mobile *string `json:"mobile,omitempty"`
	// Specifies the email address used for verification. An example value is example@email.com.
	Email *string `json:"email,omitempty"`
}

type AllowUserOpts struct {
	// Specifies whether to allow IAM users to manage access keys by themselves.
	ManageAccesskey *bool `json:"manage_accesskey,omitempty"`
	// Specifies whether to allow IAM users to change their email addresses.
	ManageEmail *bool `json:"manage_email,omitempty"`
	// Specifies whether to allow IAM users to change their mobile numbers.
	ManageMobile *bool `json:"manage_mobile,omitempty"`
	// Specifies whether to allow IAM users to change their passwords.
	ManagePassword *bool `json:"manage_password,omitempty"`
}

// UpdateProtectPolicy can modify the operation protection policy
func UpdateProtectPolicy(client *golangsdk.ServiceClient, opts *ProtectPolicyOpts, domainID string) (*ProtectPolicy, error) {
	b, err := golangsdk.BuildRequestBody(opts, "protect_policy")
	if err != nil {
		return nil, err
	}

	var resp ProtectPolicyResp
	_, err = client.Put(protectPolicyURL(client, domainID), &b, &resp, nil)
	return &resp.Body, err
}

// GetProtectPolicy retrieves details of the operation protection policy
func GetProtectPolicy(client *golangsdk.ServiceClient, domainID string) (*ProtectPolicy, error) {
	var resp ProtectPolicyResp
	_, err := client.Get(protectPolicyURL(client, domainID), &resp, nil)
	return &resp.Body, err
}
