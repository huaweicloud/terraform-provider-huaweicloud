package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateUserInformationRequestBody struct {
	User *UpdateUserInformationOption `json:"user"`
}

func (o UpdateUserInformationRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUserInformationRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateUserInformationRequestBody", string(data)}, " ")
}
