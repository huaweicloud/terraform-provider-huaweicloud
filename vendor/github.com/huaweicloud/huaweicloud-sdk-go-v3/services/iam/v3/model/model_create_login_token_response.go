package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateLoginTokenResponse struct {
	Logintoken *LoginToken `json:"logintoken,omitempty"`

	XSubjectLoginToken *string `json:"X-Subject-LoginToken,omitempty"`
	HttpStatusCode     int     `json:"-"`
}

func (o CreateLoginTokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateLoginTokenResponse struct{}"
	}

	return strings.Join([]string{"CreateLoginTokenResponse", string(data)}, " ")
}
