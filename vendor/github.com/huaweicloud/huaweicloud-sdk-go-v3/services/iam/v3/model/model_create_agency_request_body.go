package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateAgencyRequestBody struct {
	Agency *CreateAgencyOption `json:"agency"`
}

func (o CreateAgencyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgencyRequestBody struct{}"
	}

	return strings.Join([]string{"CreateAgencyRequestBody", string(data)}, " ")
}
