package auth

import (
	"github.com/chnsz/golangsdk"
)

type AuthOpts struct {
	// Application ID.
	AppId string `json:"appId" required:"true"`
	// Login account type.
	//   72: API call type.
	ClientType int `json:"clientType" required:"true"`
	// Application authentication information expiration timestamp, in seconds.
	// When the Unix timestamp of the server is greater than expireTime when the app authentication request is received,
	// the authentication fails. Example: If the application authentication information is required to expire after 10
	// minutes, expireTime = current Unix timestamp + 60*10; if The application authentication information is required
	// to never expire, expireTime = 0)
	ExpireTime *int `json:"expireTime" required:"true"`
	// A random string used to calculate application authentication information.
	// The maxLength: 64
	// The minLength: 32
	Nonce string `json:"nonce" required:"true"`
	// Enterprise ID. (When the SP application scenario is carried, if the corpId and userId fields are not carried or
	// the value is an empty string, log in as the SP default administrator)
	CorpId string `json:"corpId,omitempty"`
	// User ID. (When the userId field is not carried or the value is an empty string, log in as the enterprise default
	// administrator)
	UserId string `json:"userId,omitempty"`
	// User email.
	UserEmail string `json:"userEmail,omitempty"`
	// User name.
	UserName string `json:"userName,omitempty"`
	// User phone.
	UserPhone string `json:"userPhone,omitempty"`
	// Department code.
	DeptCode string `json:"deptCode,omitempty"`
}

// GetToken is a method to to generate a token using application information.
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
