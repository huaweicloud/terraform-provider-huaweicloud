package auth

import (
	"github.com/chnsz/golangsdk"
)

type AuthOpts struct {
	// User account (HUAWEI CLOUD meeting account).
	//   Example: zhangsan@huawei
	// Please apply for a business account in advance. For the specific application method, please refer to the
	// development process. Account has a minimum of 1 character and a maximum of 255 characters.
	Account string `json:"account" required:"true"`
	// Login client type.
	//   72: API call type.
	ClientType int `json:"clientType" required:"true"`
	// Verification code information, which is used in the verification code scenario to carry the verification code
	// information returned by the server.
	HA2 string `json:"HA2,omitempty"`
	// Whether to generate Token, the default value is 0.
	//   0: Generate token for login authentication.
	//   1: do not generate token.
	CreateTokenType *int `json:"createTokenType,omitempty"`
}

// GetToken is a method to to generate a token.
func GetToken(c *golangsdk.ServiceClient, opts AuthOpts, authorization string) (*AuthResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r AuthResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":  "application/json;charset=UTF-8",
			"Authorization": authorization,
		},
	})
	return &r, err
}

type ValidateOpts struct {
	// Verification code information, which is used in the verification code scenario to carry the verification code
	// information returned by the server.
	Token string `json:"token" required:"true"`
	// User account (HUAWEI CLOUD meeting account).
	//   Example: zhangsan@huawei
	// Please apply for a business account in advance. For the specific application method, please refer to the
	// development process. Account has a minimum of 1 character and a maximum of 255 characters.
	NeedGenNewToken bool `json:"needGenNewToken" required:"true"`
	// Login client type.
	//   72: API call type.
	NeedAccountInfo bool `json:"needAccountInfo,omitempty"`
}

// ValidateToken is a method to check whether token is available using given parameters.
func ValidateToken(c *golangsdk.ServiceClient, opts ValidateOpts) (*AuthResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	var r AuthResp
	_, err = c.Post(validateURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	})
	return &r, err
}
