package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCreateAgencyTokenResponse struct {
	Token *AgencyTokenResult `json:"token,omitempty"`

	XSubjectToken  *string `json:"X-Subject-Token,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o KeystoneCreateAgencyTokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateAgencyTokenResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCreateAgencyTokenResponse", string(data)}, " ")
}
