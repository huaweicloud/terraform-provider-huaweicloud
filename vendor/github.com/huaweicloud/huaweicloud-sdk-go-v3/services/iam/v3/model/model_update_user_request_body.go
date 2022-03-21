package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateUserRequestBody struct {
	User *UpdateUserOption `json:"user"`
}

func (o UpdateUserRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUserRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateUserRequestBody", string(data)}, " ")
}
