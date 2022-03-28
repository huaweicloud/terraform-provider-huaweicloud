package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateTemporaryAccessKeyByAgencyRequestBody struct {
	Auth *AgencyAuth `json:"auth"`
}

func (o CreateTemporaryAccessKeyByAgencyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemporaryAccessKeyByAgencyRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTemporaryAccessKeyByAgencyRequestBody", string(data)}, " ")
}
